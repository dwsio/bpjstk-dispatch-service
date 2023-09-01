package messageprocessor

import (
	"context"
	"encoding/json"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/segmentio/kafka-go"
)

func (mp *MessageProcessor) processSms(ctx context.Context, r *kafka.Reader, msg kafka.Message, consumedKafkaMsg *model.ConsumedKafkaMsg) {
	ctx, span := tracerClient.StartKafkaConsumerTracerSpan(mp.tracer, ctx, &msg.Headers, "MessageProcessor.processSms")
	defer span.End()

	publishedKafkaMsg := createPulishedKafkaMessage(consumedKafkaMsg)

	smsMsg := &model.Sms{}
	if err := json.Unmarshal(consumedKafkaMsg.Data, smsMsg); err != nil {
		mp.logKafkaMessage(ctx, false, publishedKafkaMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, publishedKafkaMsg)
		return
	}

	if err := mp.usecase.HandleSMS(ctx, publishedKafkaMsg, smsMsg); err != nil {
		mp.logKafkaMessage(ctx, false, smsMsg, err, constants.ErrorProcessingMessage)
	}

	mp.commitAndLogMsg(ctx, r, msg, smsMsg)
}
