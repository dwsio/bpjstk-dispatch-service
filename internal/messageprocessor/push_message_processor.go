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

func (mp *MessageProcessor) processPush(ctx context.Context, r *kafka.Reader, msg kafka.Message, consumedKafkaMsg *model.ConsumedKafkaMsg, producerTopic string, workerID int, traceID string) {
	ctx, span := tracerClient.StartKafkaConsumerTracerSpan(mp.tracer, ctx, tracerClient.NewKafkaHeadersCarrier(&msg.Headers), "MessageProcessor.processPush")
	defer span.End()

	publishedKafkaMsg := createPulishedKafkaMessage(consumedKafkaMsg)

	pushMsg := &model.Push{}
	if err := json.Unmarshal(consumedKafkaMsg.Data, pushMsg); err != nil {
		mp.logKafkaMessage(traceID, msg.Topic, publishedKafkaMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, publishedKafkaMsg, traceID)
		return
	}

	if err := retry.Do(func() error {
		return mp.handler.HandlePush(ctx, publishedKafkaMsg, pushMsg, traceID)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		mp.logKafkaMessage(traceID, msg.Topic, pushMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, pushMsg, traceID)
		return
	}

	mp.commitAndLogMsg(ctx, r, msg, pushMsg, traceID)
}
