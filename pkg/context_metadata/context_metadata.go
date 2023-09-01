package context_metadata

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const MetadataKey contextKey = "metadata"

type Metadata struct {
	TraceID       string
	TransactionID string
	Topic         string
	StartTime     int64
}

func SetMetadataToNewContext(ctx context.Context, traceID string, topic string) context.Context {
	metadata := Metadata{
		TraceID:       traceID,
		TransactionID: uuid.New().String(),
		Topic:         topic,
		StartTime:     time.Now().UnixMilli(),
	}

	return context.WithValue(ctx, MetadataKey, metadata)
}

func GetMetadataFromContext(ctx context.Context) (Metadata, bool) {
	metadata, ok := ctx.Value(MetadataKey).(Metadata)
	return metadata, ok
}
