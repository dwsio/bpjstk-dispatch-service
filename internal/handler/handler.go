package handler

import (
	"context"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	serviceClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/serviceclient"
	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type IHandler interface {
	HandleEmail(ctx context.Context, parentMsg *model.ConsumedKafkaMsg, emailMsg *model.Email, traceID string) error
	HandleSMS(ctx context.Context, parentMsg *model.ConsumedKafkaMsg, smsMsg *model.Sms, traceID string) error
	HandleInApp(ctx context.Context, parentMsg *model.ConsumedKafkaMsg, inappMsg *model.InApp, traceID string) error
	HandlePush(ctx context.Context, parentMsg *model.ConsumedKafkaMsg, pushMsg *model.Push, traceID string) error
}

type Handler struct {
	cfg              *config.Config
	logger           *logger.AppLogger
	tracer           trace.Tracer
	sc               serviceClient.IServiceClient
	emailWebhookUrl  string
	producerTopicMap map[string]string
	producerMap      map[string]*kafkaClient.Producer
}

func NewHandler(logger *logger.AppLogger, cfg *config.Config, tracer trace.Tracer, webhook string, producerTopicMap map[string]string, producerMap map[string]*kafkaClient.Producer) *Handler {
	return &Handler{
		logger:           logger,
		cfg:              cfg,
		sc:               serviceClient.NewServiceClient(logger, cfg, tracer),
		tracer:           tracer,
		emailWebhookUrl:  webhook,
		producerMap:      producerMap,
		producerTopicMap: producerTopicMap,
	}
}
