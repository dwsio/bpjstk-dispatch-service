package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	brokers []string
	writer  *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		brokers: brokers,
		writer:  NewWriter(brokers, topic),
	}
}

func (p *Producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.writer.WriteMessages(writeCtx, msgs...)
	if err != nil {

		return fmt.Errorf("failed to publish messages: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
