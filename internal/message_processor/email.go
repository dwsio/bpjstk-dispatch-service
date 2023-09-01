package messageprocessor

import (
	"context"
	"encoding/json"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/segmentio/kafka-go"
)

func (mp *MessageProcessor) processEmail(ctx context.Context, r *kafka.Reader, msg kafka.Message, consumedKafkaMsg *model.ConsumedKafkaMsg) {
	ctx, span := tracerClient.StartKafkaConsumerTracerSpan(mp.tracer, ctx, &msg.Headers, "MessageProcessor.processEmail")
	defer span.End()

	publishedKafkaMsg := createPulishedKafkaMessage(consumedKafkaMsg)

	emailMsg := &model.Email{}
	if err := json.Unmarshal(consumedKafkaMsg.Data, emailMsg); err != nil {
		mp.logKafkaMessage(ctx, false, publishedKafkaMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, publishedKafkaMsg)
		return
	}

	if err := mp.usecase.HandleEmail(ctx, publishedKafkaMsg, emailMsg); err != nil {
		mp.logKafkaMessage(ctx, false, emailMsg, err, constants.ErrorProcessingMessage)
	}

	mp.commitAndLogMsg(ctx, r, msg, emailMsg)
}
