package messageprocessor

import (
	"context"
	"encoding/json"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/segmentio/kafka-go"
)

func (mp *MessageProcessor) processInApp(ctx context.Context, r *kafka.Reader, msg kafka.Message, consumedKafkaMsg *model.ConsumedKafkaMsg) {
	ctx, span := tracerClient.StartKafkaConsumerTracerSpan(mp.tracer, ctx, &msg.Headers, "MessageProcessor.processInApp")
	defer span.End()

	publishedKafkaMsg := createPulishedKafkaMessage(consumedKafkaMsg)

	inappMsg := &model.InApp{}
	if err := json.Unmarshal(consumedKafkaMsg.Data, inappMsg); err != nil {
		mp.logKafkaMessage(ctx, false, publishedKafkaMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, publishedKafkaMsg)
		return
	}

	if err := mp.usecase.HandleInApp(ctx, publishedKafkaMsg, inappMsg); err != nil {
		mp.logKafkaMessage(ctx, false, inappMsg, err, constants.ErrorProcessingMessage)
	}

	mp.commitAndLogMsg(ctx, r, msg, inappMsg)
}
