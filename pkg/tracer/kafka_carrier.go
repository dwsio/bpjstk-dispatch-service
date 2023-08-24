package tracer

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// KafkaCarrier is an interface for interacting with Kafka message headers.
type KafkaCarrier interface {
	Set(key, value string)
	Get(key string) string
	Keys() []string
}

// StartKafkaConsumerTracerSpan starts a new span for a Kafka consumer operation.
func StartKafkaConsumerTracerSpan(tracer trace.Tracer, consumerCtx context.Context, carrier KafkaCarrier, operationName string) (context.Context, trace.Span) {
	propagatedCtx := otel.GetTextMapPropagator().Extract(consumerCtx, carrier)

	consumerCtx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(consumerCtx, trace.SpanContextFromContext(propagatedCtx)),
		operationName,
	)

	return consumerCtx, span
}

// InjectKafkaTracingHeadersToCarrier injects tracing headers into the carrier.
func InjectKafkaTracingHeadersToCarrier(producerCtx context.Context, carrier KafkaCarrier) {
	otel.GetTextMapPropagator().Inject(producerCtx, carrier)
}

// KafkaHeadersCarrier implements the KafkaCarrier interface using kafka.Header.
type KafkaHeadersCarrier struct {
	Headers *[]kafka.Header
}

// Set sets a key-value pair in the Kafka headers.
func (c KafkaHeadersCarrier) Set(key, value string) {
	*c.Headers = append(*c.Headers, kafka.Header{Key: key, Value: []byte(value)})
}

// Get retrieves the value for a key from the Kafka headers.
func (c KafkaHeadersCarrier) Get(key string) string {
	for _, h := range *c.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

// Keys returns a list of keys present in the Kafka headers.
func (c KafkaHeadersCarrier) Keys() []string {
	keys := make([]string, len(*c.Headers))
	for i, h := range *c.Headers {
		keys[i] = h.Key
	}
	return keys
}

// NewKafkaHeadersCarrier creates a new instance of KafkaHeadersCarrier.
func NewKafkaHeadersCarrier(headers *[]kafka.Header) KafkaHeadersCarrier {
	return KafkaHeadersCarrier{Headers: headers}
}
