package messageprocessor

import (
	"context"
	"encoding/json"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/segmentio/kafka-go"
)

func (mp *MessageProcessor) processPush(ctx context.Context, r *kafka.Reader, msg kafka.Message, consumedKafkaMsg *model.ConsumedKafkaMsg) {
	ctx, span := tracerClient.StartKafkaConsumerTracerSpan(mp.tracer, ctx, &msg.Headers, "MessageProcessor.processPush")
	defer span.End()

	publishedKafkaMsg := createPulishedKafkaMessage(consumedKafkaMsg)

	pushMsg := &model.Push{}
	if err := json.Unmarshal(consumedKafkaMsg.Data, pushMsg); err != nil {
		mp.logKafkaMessage(ctx, false, publishedKafkaMsg, err, constants.ErrorProcessingMessage)
		mp.commitAndLogMsg(ctx, r, msg, publishedKafkaMsg)
		return
	}

	if err := mp.usecase.HandlePush(ctx, publishedKafkaMsg, pushMsg); err != nil {
		mp.logKafkaMessage(ctx, false, pushMsg, err, constants.ErrorProcessingMessage)
	}

	mp.commitAndLogMsg(ctx, r, msg, pushMsg)
}
