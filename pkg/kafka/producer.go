package kafka

import (
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"

	"log"
)

var (
	producer sarama.AsyncProducer
	brokers  []string
)

func InitProducer(cfg *models.KafkaConfig) string {
	brokers = cfg.Brokers
	if len(brokers) == 0 {
		return errcode.Msg(errcode.ErrForKafka)
	}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Compression = sarama.CompressionSnappy
	saramaConfig.Producer.Retry.Max = 3
	var err error
	producer, err = sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		return errcode.Msg(errcode.ERROR)
	}

	go func() {
		for err = range producer.Errors() {
			log.Printf("kafka producer error: %v\n", err)
		}
	}()

	go func() {
		for success := range producer.Successes() {
			log.Printf("Message published to topic=%s,partition=%d,offset=%d\n",
				success.Topic, success.Partition, success.Offset)
		}
	}()

	log.Printf("kafka producer initialized with brokers:%v\n", brokers)
	return errcode.Msg(errcode.SUCCESS)
}

func Close() {
	if producer != nil {
		if err := producer.Close(); err != nil {
			log.Printf("close kafka producer failed,err:%v\n", err)
		}
	}
}

func Publish(topic string, key string, value interface{}) string {
	if producer == nil {
		return errcode.Msg(errcode.ErrForKafka)
	}
	data, err := json.Marshal(value)
	if err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	select {
	case producer.Input() <- msg:
		return errcode.Msg(errcode.SUCCESS)
	case <-time.After(5 * time.Second):
		return errcode.Msg(errcode.KafTimeOut)
	}
}
