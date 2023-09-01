package messageprocessor

import (
	"context"
	"encoding/json"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	serviceMetrics "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/service_metrics"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/usecase"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	contextMd "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/context_metadata"
	loggerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type MessageProcessor struct {
	cfg              *config.Config
	logger           *loggerClient.AppLogger
	usecase          *usecase.Usecase
	tracer           trace.Tracer
	producerTopicMap map[string]string
	serviceMetrics   *serviceMetrics.ServiceMetrics
}

func NewMessageProcessor(logger *loggerClient.AppLogger, cfg *config.Config, usecase *usecase.Usecase, tracer trace.Tracer, producerTopicMap map[string]string, serviceMetrics *serviceMetrics.ServiceMetrics) *MessageProcessor {
	return &MessageProcessor{
		logger:           logger,
		cfg:              cfg,
		usecase:          usecase,
		tracer:           tracer,
		producerTopicMap: producerTopicMap,
		serviceMetrics:   serviceMetrics,
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
			mp.serviceMetrics.ErrorKafkaConsume.Add(ctx, 1, metric.WithAttributes())
			mp.logKafkaMessage(ctx, false, nil, err, "Failed to fetch kafka message")
			continue
		}
		mp.serviceMetrics.SuccessKafkaConsume.Add(ctx, 1, metric.WithAttributes())

		traceID := getValueFromKafkaHeaders(createKafkaHeadersMap(fetchedMessage.Headers), "trace_id")
		ctx = contextMd.SetMetadataToNewContext(ctx, traceID, fetchedMessage.Topic)

		consumedKafkaMsg := &model.ConsumedKafkaMsg{}
		if err := json.Unmarshal(fetchedMessage.Value, consumedKafkaMsg); err != nil {
			mp.logKafkaMessage(ctx, false, nil, err, constants.ErrorProcessingMessage)
			mp.commitAndLogMsg(ctx, r, fetchedMessage, "")
			continue
		}

		mp.logKafkaMessage(ctx, true, consumedKafkaMsg, nil, "Kafka message received and is being processed")

		if processorFunc, exists := mp.getProcessorFunc(consumedKafkaMsg.CategoryName); exists {
			processorFunc(ctx, r, fetchedMessage, consumedKafkaMsg)
		} else {
			mp.logKafkaMessage(ctx, false, consumedKafkaMsg, nil, "Unsupported message category")
			mp.commitAndLogMsg(ctx, r, fetchedMessage, consumedKafkaMsg)
		}
	}
}
