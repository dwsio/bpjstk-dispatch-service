package messageprocessor

import (
	"context"
	"encoding/json"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/handler"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/utils"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"

	"github.com/avast/retry-go"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

const (
	retryAttempts = 1
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

type MessageProcessor struct {
	cfg              *config.Config
	logger           *logger.AppLogger
	handler          *handler.Handler
	tracer           trace.Tracer
	producerTopicMap map[string]string
}

func NewMessageProcessor(logger *logger.AppLogger, cfg *config.Config, handler *handler.Handler, tracer trace.Tracer, producerTopicMap map[string]string) *MessageProcessor {
	return &MessageProcessor{
		logger:           logger,
		cfg:              cfg,
		handler:          handler,
		tracer:           tracer,
		producerTopicMap: producerTopicMap,
	}
}

func (mp *MessageProcessor) ProcessMessage(ctx context.Context, r *kafka.Reader, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fetchedMessage, err := r.FetchMessage(ctx)
		if err != nil {
			mp.logKafkaMessage("", r.Config().Topic, "", nil, "Failed to fetch kafka message")
			continue
		}

		traceID := getValueFromKafkaHeaders(createKafkaHeadersMap(fetchedMessage.Headers), "trace_id")

		consumedKafkaMsg := &model.ConsumedKafkaMsg{}
		if err := json.Unmarshal(fetchedMessage.Value, consumedKafkaMsg); err != nil {
			mp.logKafkaMessage(traceID, fetchedMessage.Topic, "", err, constants.ErrorProcessingMessage)
			mp.commitAndLogMsg(ctx, r, fetchedMessage, "", traceID)
			continue
		}

		dataLog, err := utils.StructToString(consumedKafkaMsg)
		if err != nil {
			continue
		}

		mp.logKafkaMessage(traceID, fetchedMessage.Topic, dataLog, nil, "Kafka message received and is being processed")

		if processorFunc, exists := mp.getProcessorFunc(consumedKafkaMsg.CategoryName); exists {
			processorFunc(ctx, r, fetchedMessage, consumedKafkaMsg, mp.producerTopicMap[consumedKafkaMsg.CategoryName], workerID, traceID)
		} else {
			mp.logKafkaMessage(traceID, fetchedMessage.Topic, dataLog, nil, "Unsupported message category")
			mp.commitAndLogMsg(ctx, r, fetchedMessage, dataLog, traceID)
		}
	}
}
