package kfk

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaReader struct {
	r *kafka.Reader
}

func GetReader(brokers []string, groupId, topic string) *KafkaReader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	kr := &KafkaReader{
		r: r,
	}
	go kr.readMsg()
	return kr
}

func (kr *KafkaReader) readMsg() {
	for {
		m, err := kr.r.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("kafka receiver read msg err", zap.Error(err))
			continue
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
