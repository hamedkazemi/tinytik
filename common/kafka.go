package common

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

type KafkaConnection struct {
	KW *kafka.Writer
}

var KafkaInstance *KafkaConnection
var onceKf sync.Once

func KafkaCon() *KafkaConnection {
	onceKf.Do(func() {
		fmt.Println(Config.Kafka.Ip + ":" + Config.Kafka.Port)
		fmt.Println(Config.Kafka.Topic)
		w := kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{Config.Kafka.Ip + ":" + Config.Kafka.Port},
			Topic:        Config.Kafka.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
		})
		KafkaInstance = &KafkaConnection{KW: w}
	})
	return KafkaInstance
}

func PublishMessageToKafka(key string, message string) error {
	err := KafkaCon().KW.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(message),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
