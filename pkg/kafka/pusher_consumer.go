package kafka

import (
	"IM_chat/pkg/errcode"
	"context"
	"github.com/IBM/sarama"
	"log"
	"strconv"
	"sync"
)

type PersistConsumer struct {
	reader  sarama.ConsumerGroup
	topics  []string
	handler *persistHandler
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// persistHandler 实现了 sarama.ConsumerGroupHandler 接口。
// 它把每条消息的具体业务处理委托给 handleFunc，便于外部注入业务逻辑。
type persistHandler struct {
	ready      chan bool
	handleFunc func(topic string, key []byte, value []byte) error
}

// NewPersistConsumer 创建一个 Kafka 消费组消费者。
// group_id 是消费组标识；每消费到一条消息，都会调用一次 handleFunc。
func NewPersistConsumer(
	brokers []string,
	group_id int64,
	topics []string,
	handleFunc func(topic string, key []byte, value []byte) error,
) (*PersistConsumer, string) {
	if len(brokers) == 0 || len(topics) == 0 || handleFunc == nil {
		return nil, errcode.Msg(errcode.ErrForKafka)
	}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Version = sarama.V2_1_0_0
	groupID := strconv.Itoa(int(group_id))
	reader, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		return nil, errcode.Msg(errcode.ERROR)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &PersistConsumer{
		reader: reader,
		topics: topics,
		handler: &persistHandler{
			ready:      make(chan bool),
			handleFunc: handleFunc,
		},
		ctx:    ctx,
		cancel: cancel,
	}, errcode.Msg(errcode.SUCCESS)
}

// Start 在后台 goroutine 中启动消费循环。
// 发生 rebalance 后 Consume 可能返回，循环会继续拉起，直到调用 Stop。
func (c *PersistConsumer) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.reader.Consume(c.ctx, c.topics, c.handler); err != nil {
				log.Println("PersistConsumer error:", err)
			}
			if c.ctx.Err() != nil {
				return
			}
			c.handler.ready = make(chan bool)
		}
	}()
	<-c.handler.ready
	log.Println("persistConsumer started")
}

// Stop 取消消费上下文，等待消费循环退出，并关闭 consumer group 客户端。
func (c *PersistConsumer) Stop() error {
	c.cancel()
	c.wg.Wait()
	return c.reader.Close()
}

// Setup 在每次消费组会话开始时被调用。
func (h *persistHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

// Cleanup 在会话结束时被调用（例如发生 rebalance）。
func (h *persistHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 负责消费当前 claim（分区分配）中的消息。
// 只有 handleFunc 成功后才调用 MarkMessage，表示这条消息“已处理完成”，
// 这样提交的 offset 语义才是可靠的“处理完成再提交”。
func (h *persistHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var key []byte
		if msg.Key != nil {
			key = append([]byte(nil), msg.Key...)
		}

		err := h.handleFunc(msg.Topic, key, msg.Value)
		if err != nil {
			log.Printf("persist consume failed topic=%s partition=%d offset=%d err=%v",
				msg.Topic, msg.Partition, msg.Offset, err)
			continue
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
