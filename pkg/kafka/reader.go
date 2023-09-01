package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	minBytes                     = 10e3
	maxBytes                     = 10e6
	dialTimeout                  = 3 * time.Minute
	readerHeartbeatInterval      = 3 * time.Second
	readerCommitInterval         = 0
	readerPartitionWatchInterval = 5 * time.Second
	readerMaxAttempts            = 3
	readerMaxWait                = 1 * time.Second
	readerQueueCapacity          = 100
	readerReadBatchTimeout       = 1 * time.Second
	readerReadLagInterval        = -1
)

func NewReader(brokers []string, groupID string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                brokers,
		GroupID:                groupID,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          readerQueueCapacity,
		ReadBatchTimeout:       readerReadBatchTimeout,
		HeartbeatInterval:      readerHeartbeatInterval,
		ReadLagInterval:        readerReadLagInterval,
		CommitInterval:         readerCommitInterval,
		PartitionWatchInterval: readerPartitionWatchInterval,
		MaxAttempts:            readerMaxAttempts,
		MaxWait:                readerMaxWait,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
	})
}
