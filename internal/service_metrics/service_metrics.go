package service_metrics

import (
	"go.opentelemetry.io/otel/metric"
)

type ServiceMetrics struct {
	SuccessKafkaConsume metric.Int64Counter
	ErrorKafkaConsume   metric.Int64Counter
	SuccessKafkaPublish metric.Int64Counter
	ErrorKafkaPublish   metric.Int64Counter
}

func NewServiceMetrics(meter metric.Meter) *ServiceMetrics {
	successKafkaConsume, _ := meter.Int64Counter(
		"success_kafka_consume",
		metric.WithDescription("The total number of success kafka consume"),
	)

	errorKafkaConsume, _ := meter.Int64Counter(
		"error_kafka_consume",
		metric.WithDescription("The total number of error kafka consume"),
	)

	successKafkaPublish, _ := meter.Int64Counter(
		"success_kafka_publish",
		metric.WithDescription("The total number of success kafka publish"),
	)

	errorKafkaPublish, _ := meter.Int64Counter(
		"error_kafka_publish",
		metric.WithDescription("The total number of error kafka pulish"),
	)

	return &ServiceMetrics{
		SuccessKafkaConsume: successKafkaConsume,
		ErrorKafkaConsume:   errorKafkaConsume,
		SuccessKafkaPublish: successKafkaPublish,
		ErrorKafkaPublish:   errorKafkaPublish,
	}
}
