package kfk

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type LogData struct {
	Topic string
	Data  []byte
}

type KafkaWriter struct {
	w    *kafka.Writer
	data chan LogData
}

func GetWriter(addr string) *KafkaWriter {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Balancer: &kafka.LeastBytes{},
	}
	kw := KafkaWriter{
		w:    writer,
		data: make(chan LogData, 100),
	}
	go kw.sendMsg()
	return &kw
}

func (kw *KafkaWriter) Close() {
	if kw.w != nil {
		_ = kw.w.Close()
	}
}

func (kw *KafkaWriter) Send(data LogData) {
	kw.data <- data
}

func (kw *KafkaWriter) sendMsg() {
	for {
		select {
		case data := <-kw.data:
			msg := kafka.Message{
				Topic: data.Topic,
				Key:   []byte("logMessage"),
				Value: data.Data,
			}
			var err error
			const retries = 3
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			//defer cancel()
			for i := 0; i < retries; i++ {

				if err = kw.w.WriteMessages(ctx, msg); err == nil {
					break
				}
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					continue
				}
				if err != nil {
					zap.L().Error("kafka send log writer msg err", zap.Error(err))
				}
			}
			cancel()
		}
	}
}
