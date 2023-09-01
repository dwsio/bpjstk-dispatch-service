package tracer

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type KafkaCarrier interface {
	Set(key, value string)
	Get(key string) string
	Keys() []string
}

type KafkaHeadersCarrier struct {
	Headers *[]kafka.Header
}

func (c KafkaHeadersCarrier) Set(key, value string) {
	*c.Headers = append(*c.Headers, kafka.Header{Key: key, Value: []byte(value)})
}

func (c KafkaHeadersCarrier) Get(key string) string {
	for _, h := range *c.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c KafkaHeadersCarrier) Keys() []string {
	keys := make([]string, len(*c.Headers))
	for i, h := range *c.Headers {
		keys[i] = h.Key
	}
	return keys
}

func StartKafkaConsumerTracerSpan(tracer trace.Tracer, consumerCtx context.Context, headers *[]kafka.Header, operationName string) (context.Context, trace.Span) {
	propagatedCtx := otel.GetTextMapPropagator().Extract(consumerCtx, KafkaHeadersCarrier{Headers: headers})

	consumerCtx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(consumerCtx, trace.SpanContextFromContext(propagatedCtx)),
		operationName,
	)

	return consumerCtx, span
}

func InjectKafkaTracingHeadersToCarrier(producerCtx context.Context, headers *[]kafka.Header) {
	otel.GetTextMapPropagator().Inject(producerCtx, KafkaHeadersCarrier{Headers: headers})
}
