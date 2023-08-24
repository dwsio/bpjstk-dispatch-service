package messageprocessor

import (
	"context"
	"encoding/json"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/avast/retry-go"
	"github.com/segmentio/kafka-go"
)

func (mp *MessageProcessor) processInApp(ctx context.Context, r *kafka.Reader, msg kafka.Message, consumedKafkaMsg *model.ConsumedKafkaMsg, producerTopic string, workerID int, traceID string) {
	ctx, span := tracerClient.StartKafkaConsumerTracerSpan(mp.tracer, ctx, tracerClient.NewKafkaHeadersCarrier(&msg.Headers), "MessageProcessor.processInApp")
	defer span.End()

	publishedKafkaMsg := createPulishedKafkaMessage(consumedKafkaMsg)

	inappMsg := &model.InApp{}
	if err := json.Unmarshal(consumedKafkaMsg.Data, inappMsg); err != nil {
		mp.logKafkaMessage(traceID, msg.Topic, publishedKafkaMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, publishedKafkaMsg, traceID)
		return
	}

	if err := retry.Do(func() error {
		return mp.handler.HandleInApp(ctx, publishedKafkaMsg, inappMsg, traceID)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		mp.logKafkaMessage(traceID, msg.Topic, inappMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, inappMsg, traceID)
		return
	}

	mp.commitAndLogMsg(ctx, r, msg, inappMsg, traceID)
}
