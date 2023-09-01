package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	messageProcessor "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/message_processor"
	serviceMetrics "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/service_metrics"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/usecase"
	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	loggerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	metricClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/metric"
	metricServer "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/metric_server"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
)

type Server struct {
	cfg              *config.Config
	appLogger        *loggerClient.AppLogger
	appTracer        *tracerClient.AppTracer
	appMetric        *metricClient.AppMetric
	usecase          *usecase.Usecase
	producerTopicMap map[string]string
	producerMap      map[string]*kafkaClient.Producer
	serviceMetrics   *serviceMetrics.ServiceMetrics
	metricServer     *metricServer.MetricServer
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run(channel string, priority string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := s.setupLogger(); err != nil {
		return fmt.Errorf("logger setup failed: %w", err)
	}

	if err := s.setupTracer(ctx); err != nil {
		return fmt.Errorf("tracer setup failed: %w", err)
	}
	defer func() {
		if err := s.appTracer.TraceProvider.Shutdown(ctx); err != nil {
			fmt.Println("Failed to shutdown tracer:", err)
		}
	}()

	if err := s.setupMetric(ctx, cancel); err != nil {
		return fmt.Errorf("metric setup failed: %w", err)
	}
	defer func() {
		if err := s.metricServer.Shutdown(ctx); err != nil {
			fmt.Println("Failed to shutdown metric server:", err)
		}
	}()

	if err := s.setupKafka(ctx); err != nil {
		return fmt.Errorf("kafka setup failed: %w", err)
	}

	producerTopics := strings.Split(s.cfg.KafkaTopic.Producer, ",")
	producerBrokers := strings.Split(s.cfg.Kafka.ProducerBrokers, ",")
	s.producerMap = s.createProducerMap(producerTopics, producerBrokers)
	defer func() {
		for _, producer := range s.producerMap {
			if closeErr := producer.Close(); closeErr != nil {
				fmt.Println("Failed to close Kafka producer:", closeErr)
			}
		}
	}()

	s.setupUsecase()

	messageProcessor := s.createMessageProcessor()
	consumerTopics := s.prepareConsumerTopics(channel, priority)
	if err := s.startConsumers(ctx, consumerTopics, messageProcessor); err != nil {
		fmt.Println("Consumer start failed:", err)
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) setupLogger() error {
	var err error

	s.appLogger = loggerClient.NewAppLogger()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) setupTracer(ctx context.Context) error {
	var err error

	s.appTracer, err = tracerClient.NewAppTracer(ctx, s.cfg.Tracer, s.cfg.Project.ServiceName, s.cfg.Project.Version)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) setupMetric(ctx context.Context, cancel context.CancelFunc) error {
	var err error

	s.appMetric, err = metricClient.NewAppMetric(ctx, s.cfg.Metric, s.cfg.Project.ServiceName, s.cfg.Project.Version)
	if err != nil {
		return err
	}

	s.serviceMetrics = serviceMetrics.NewServiceMetrics(s.appMetric.Meter)
	s.metricServer = metricServer.NewMetricServer(s.cfg.Metric.Path)

	go func() {
		defer cancel()

		if err := s.metricServer.Run(s.cfg.Metric.Port); err != nil {
			fmt.Println("Failed to run metric server:", err)
		}
	}()

	return nil
}

func (s *Server) setupKafka(ctx context.Context) error {
	producerBrokers := strings.Split(s.cfg.Kafka.ProducerBrokers, ",")
	fmt.Println("Producers broker:", producerBrokers)

	producerTopics := strings.Split(s.cfg.KafkaTopic.Producer, ",")
	s.producerTopicMap = s.createTopicMap(producerTopics)

	if err := s.initKafkaTopics(ctx, producerBrokers[0], producerTopics); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupUsecase() {
	s.usecase = usecase.NewUsecase(s.appLogger, s.cfg, s.appTracer.Tracer, s.cfg.ProviderClient.EmailProvider.Webhook, s.producerTopicMap, s.producerMap, s.serviceMetrics)
}

func (s *Server) createMessageProcessor() *messageProcessor.MessageProcessor {
	return messageProcessor.NewMessageProcessor(s.appLogger, s.cfg, s.usecase, s.appTracer.Tracer, s.producerTopicMap, s.serviceMetrics)
}
