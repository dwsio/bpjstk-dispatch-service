package kafka

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	ProducerBrokers string `mapstructure:"PRODUCER_BROKERS"`
	ConsumerBrokers string `mapstructure:"CONSUMER_BROKERS"`
	GroupID         string `mapstructure:"GROUP_ID"`
	PoolSize        string `mapstructure:"POOL_SIZE"`
	Partition       string `mapstructure:"PARTITION"`
}

const (
	ReplicationFactor = 1
)

func NewKafkaConn(ctx context.Context, kafkaCfg *Config) (*kafka.Conn, error) {
	producerBrokers := strings.Split(kafkaCfg.ProducerBrokers, ",")

	return kafka.DialContext(ctx, "tcp", producerBrokers[0])
}
