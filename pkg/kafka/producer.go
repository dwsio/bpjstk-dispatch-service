package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

// Producer represents a Kafka message producer.
type Producer struct {
	brokers []string
	writer  *kafka.Writer
}

// NewProducer creates a new Kafka producer.
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		brokers: brokers,
		writer:  NewWriter(brokers, topic),
	}
}

// PublishMessage publishes Kafka messages.
func (p *Producer) PublishMessage(ctx context.Context, tracer trace.Tracer, msgs ...kafka.Message) error {
	// Start a new span for this publish operation.
	ctx, span := tracer.Start(ctx, "Producer.PublishMessage")
	defer span.End()

	// Use a new context for WriteMessages with timeout.
	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Write messages to Kafka.
	err := p.writer.WriteMessages(writeCtx, msgs...)
	if err != nil {
		// Handle the error appropriately.
		return fmt.Errorf("failed to publish messages: %w", err)
	}

	return nil
}

// Close closes the Kafka producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
