package kafka

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
)

// Worker defines the type for the worker function that processes Kafka messages.
type Worker func(ctx context.Context, r *kafka.Reader, workerID int)

// Consumer provides methods to manage Kafka consumers and their workers.
type Consumer struct {
	brokers []string // List of Kafka brokers
}

// NewConsumer creates a new instance of the Kafka Consumer.
func NewConsumer(brokers []string) *Consumer {
	return &Consumer{
		brokers: brokers,
	}
}

// StartWorkers starts a pool of worker goroutines to process Kafka messages.
func (c *Consumer) StartWorkers(ctx context.Context, groupID string, consumerTopic string, poolSize int, worker Worker) {
	var wg sync.WaitGroup

	// Launch poolSize number of worker goroutines
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Create a new Kafka reader for this worker
			reader := NewReader(c.brokers, groupID, consumerTopic)
			defer reader.Close()

			// Execute the worker function for processing messages
			worker(ctx, reader, workerID)
		}(i)
	}
	wg.Wait()
}
