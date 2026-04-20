package kafka

import (
	"IM_chat/pkg/errcode"
	"context"
	"github.com/IBM/sarama"
	"log"
	"strconv"
	"sync"
	"time"
)

type PersistConsumer struct {
	reader  sarama.ConsumerGroup
	topics  []string
	handler *persistHandler
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

type persistHandler struct {
	ready      chan bool
	handleFunc func(topic string, key []byte, value []byte) error
}

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

			select {
			case <-c.ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}()

	select {
	case <-c.handler.ready:
		log.Println("persistConsumer started")
	case <-time.After(10 * time.Second):
		log.Println("persistConsumer: start timed out, will keep retrying in background")
	}
}

func (c *PersistConsumer) Stop() error {
	c.cancel()
	c.wg.Wait()
	return c.reader.Close()
}

func (h *persistHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *persistHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

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
