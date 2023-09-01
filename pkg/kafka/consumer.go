package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Worker func(ctx context.Context, r *kafka.Reader, workerID int)

type Consumer struct {
	brokers []string
}

func NewConsumer(brokers []string) *Consumer {
	return &Consumer{
		brokers: brokers,
	}
}

func (c *Consumer) StartWorkers(ctx context.Context, groupID string, consumerTopic string, poolSize int, worker Worker) {
	var wg sync.WaitGroup

	for i := 0; i < poolSize; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			reader := NewReader(c.brokers, groupID, consumerTopic)

			defer func() {
				if err := reader.Close(); err != nil {
					fmt.Println("Failed to close reader:", err)
				}
			}()

			worker(ctx, reader, workerID)
		}(i)
	}

	wg.Wait()
}
