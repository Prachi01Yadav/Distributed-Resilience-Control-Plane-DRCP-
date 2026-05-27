package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
	logger *zap.Logger
}

func NewConsumer(brokers []string, topic, groupID string, logger *zap.Logger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader: r,
		logger: logger,
	}
}

func (c *Consumer) ReadMessages(ctx context.Context, handler func(msg []byte) error) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.logger.Error("Failed to read message", zap.Error(err))
			if ctx.Err() != nil {
				return
			}
			continue
		}

		if err := handler(m.Value); err != nil {
			c.logger.Error("Handler error", zap.Error(err))
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
