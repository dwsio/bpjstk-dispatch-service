package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

const (
	writerBatchSize    = 100
	writerBatchBytes   = 1048576
	writerReadTimeout  = 10 * time.Second
	writerWriteTimeout = 10 * time.Second
	writerBatchTimeout = 1 * time.Second
	writerRequiredAcks = -1
	writerMaxAttempts  = 3
	writerAsync        = false
)

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
