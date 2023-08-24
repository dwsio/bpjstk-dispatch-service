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

func (mp *MessageProcessor) processEmail(ctx context.Context, r *kafka.Reader, msg kafka.Message, consumedKafkaMsg *model.ConsumedKafkaMsg, producerTopic string, workerID int, traceID string) {
	ctx, span := tracerClient.StartKafkaConsumerTracerSpan(mp.tracer, ctx, tracerClient.NewKafkaHeadersCarrier(&msg.Headers), "MessageProcessor.processEmail")
	defer span.End()

	publishedKafkaMsg := createPulishedKafkaMessage(consumedKafkaMsg)

	emailMsg := &model.Email{}
	if err := json.Unmarshal(consumedKafkaMsg.Data, emailMsg); err != nil {
		mp.logKafkaMessage(traceID, msg.Topic, publishedKafkaMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, publishedKafkaMsg, traceID)
		return
	}

	if err := retry.Do(func() error {
		return mp.handler.HandleEmail(ctx, publishedKafkaMsg, emailMsg, traceID)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		mp.logKafkaMessage(traceID, msg.Topic, emailMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, emailMsg, traceID)
		return
	}

	mp.commitAndLogMsg(ctx, r, msg, emailMsg, traceID)
}
