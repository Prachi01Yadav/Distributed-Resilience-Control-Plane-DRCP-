package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewProducer(brokers []string, topic string, logger *zap.Logger) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		writer: w,
		logger: logger,
	}
}

func (p *Producer) ProduceTelemetry(ctx context.Context, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
	if err != nil {
		p.logger.Error("Failed to write to kafka", zap.Error(err))
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
