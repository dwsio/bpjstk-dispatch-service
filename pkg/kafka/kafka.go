package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type Config struct {
	ProducerBroker string
	ConsumerBroker string
	GroupID        string
	PoolSize       string
	Partition      string
}

const (
	// Topic config
	ReplicationFactor = 1

	// Reader config
	minBytes                     = 10e3 // 10KB
	maxBytes                     = 10e6 // 10MB
	dialTimeout                  = 3 * time.Minute
	readerHeartbeatInterval      = 3 * time.Second
	readerCommitInterval         = 0
	readerPartitionWatchInterval = 5 * time.Second
	readerMaxAttempts            = 3
	readerMaxWait                = 1 * time.Second
	readerQueueCapacity          = 100
	readerReadBatchTimeout       = 1 * time.Second
	readerReadLagInterval        = -1

	// Writer config
	writerBatchSize    = 100
	writerBatchBytes   = 1048576
	writerReadTimeout  = 10 * time.Second
	writerWriteTimeout = 10 * time.Second
	writerBatchTimeout = 1 * time.Second
	writerRequiredAcks = -1
	writerMaxAttempts  = 3
	writerAsync        = false
)

func NewKafkaConn(ctx context.Context, broker string) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", broker)
}

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

func NewWriter(brokers []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
		Balancer:         &kafka.LeastBytes{},
		MaxAttempts:      writerMaxAttempts,
		BatchSize:        writerBatchSize,
		BatchBytes:       writerBatchBytes,
		BatchTimeout:     writerBatchTimeout,
		RequiredAcks:     writerRequiredAcks,
		ReadTimeout:      writerReadTimeout,
		WriteTimeout:     writerWriteTimeout,
		CompressionCodec: &compress.SnappyCodec,
		Logger:           nil,
		ErrorLogger:      nil,
		Async:            writerAsync,
	})
}
