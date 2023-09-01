package usecase

import (
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	providerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/provider_client"
	serviceMetrics "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/service_metrics"
	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	loggerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type Usecase struct {
	cfg              *config.Config
	logger           *loggerClient.AppLogger
	tracer           trace.Tracer
	sc               providerClient.IProviderClient
	emailWebhookUrl  string
	producerTopicMap map[string]string
	producerMap      map[string]*kafkaClient.Producer
	serviceMetrics   *serviceMetrics.ServiceMetrics
}

func NewUsecase(logger *loggerClient.AppLogger, cfg *config.Config, tracer trace.Tracer, webhook string, producerTopicMap map[string]string, producerMap map[string]*kafkaClient.Producer, serviceMetrics *serviceMetrics.ServiceMetrics) *Usecase {
	return &Usecase{
		logger:           logger,
		cfg:              cfg,
		sc:               providerClient.NewProviderClient(logger, cfg, tracer),
		tracer:           tracer,
		emailWebhookUrl:  webhook,
		producerMap:      producerMap,
		producerTopicMap: producerTopicMap,
		serviceMetrics:   serviceMetrics,
	}
}
