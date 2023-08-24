package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	// Import necessary packages ...

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/handler"
	mp "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/messageprocessor"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	loggerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/segmentio/kafka-go"
)

// Server represents the main application structure.
type Server struct {
	cfg              *config.Config                   // Configuration instance for application settings
	logger           *loggerClient.AppLogger          // Logger instance for managing application logs
	tracer           *tracerClient.AppTracer          // Tracer instance for distributed tracing
	producerTopicMap map[string]string                // Mapping of producer topics
	producerMap      map[string]*kafkaClient.Producer // Mapping of Kafka producer instances
	handler          *handler.Handler                 // Handler instance for processing messages
}

// NewServer creates a new instance of the Server.
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// Run is the entry point for running the server.
func (s *Server) Run(channel string, priority string) error {
	// Set up a context for graceful shutdown ...
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Initialize the logger ...
	if err := s.setupLogger(); err != nil {
		return fmt.Errorf("logger setup failed: %w", err)
	}

	// Set up tracing for distributed tracing functionality ...
	if err := s.setupTracer(ctx); err != nil {
		return fmt.Errorf("tracer setup failed: %w", err)
	}
	defer func() {
		// Perform tracer shutdown ...
		if shutdownErr := s.tracer.TraceProvider.Shutdown(ctx); shutdownErr != nil {
			fmt.Println("Failed to shutdown tracer:", shutdownErr)
		}
	}()

	// Set up Kafka connection and topics ...
	if err := s.setupKafka(ctx); err != nil {
		return fmt.Errorf("kafka setup failed: %w", err)
	}

	// Parse Kafka broker and topic configuration ...
	producerBrokers := strings.Split(s.cfg.Kafka.ProducerBroker, ",")
	s.producerMap = s.createProducerMap(strings.Split(s.cfg.KafkaTopic.Producer, ","), producerBrokers)
	defer func() {
		// Close Kafka producers on exit ...
		for _, producer := range s.producerMap {
			if closeErr := producer.Close(); closeErr != nil {
				fmt.Println("Failed to close Kafka producer:", closeErr)
			}
		}
	}()

	// Set up the message handler ...
	s.setupHandler()

	// Create a message processor and prepare consumer topics based on channel and priority ...
	messageProcessor := s.createMessageProcessor()
	consumerTopics := s.prepareConsumerTopics(channel, priority)

	if err := s.startConsumers(ctx, consumerTopics, messageProcessor); err != nil {
		fmt.Println("Consumer start failed:", err)
		return err
	}

	// Wait for the context to be done, indicating a shutdown signal ...
	<-ctx.Done()
	return nil
}

// setupLogger initializes the logger for the application.
func (s *Server) setupLogger() error {
	var err error
	// Initialize the logger using the provided encoding method ...
	s.logger, err = loggerClient.NewAppLogger(s.cfg.Logger)
	if err != nil {
		return err
	}

	return nil
}

// setupTracer initializes the OpenTelemetry tracer for distributed tracing.
func (s *Server) setupTracer(ctx context.Context) error {
	var err error
	// Initialize the tracer with the provided settings ...
	s.tracer, err = tracerClient.NewAppTracer(ctx, s.cfg.Tracer, s.cfg.Project.ServiceName, s.cfg.Project.Version)
	if err != nil {
		return err
	}

	return nil
}

// setupKafka initializes the Kafka connection and topics for message communication.
func (s *Server) setupKafka(ctx context.Context) error {
	// Parse Kafka broker and topic configuration ...
	producerBrokers := strings.Split(s.cfg.Kafka.ProducerBroker, ",")
	fmt.Println("Producers broker:", producerBrokers)

	// Create a mapping of producer topics for easier access ...
	producerTopics := strings.Split(s.cfg.KafkaTopic.Producer, ",")
	s.producerTopicMap = s.createTopicMap(producerTopics)

	// Initialize Kafka topics and connections ...
	if err := s.initKafkaTopics(ctx, producerBrokers[0], producerTopics); err != nil {
		return err
	}
	return nil
}

// initKafkaTopics initializes Kafka topics.
func (s *Server) initKafkaTopics(ctx context.Context, producerBroker string, topics []string) error {
	// Establish a connection with the Kafka controller ...
	conn, err := kafka.DialContext(ctx, "tcp", producerBroker)
	if err != nil {
		return err
	}
	// Close the Kafka connection on exit ...
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Println("Failed to close Kafka connection:", closeErr)
		}
	}()

	fmt.Println("Established new kafka controller connection:", producerBroker)

	partition, err := strconv.Atoi(s.cfg.Kafka.Partition)
	if err != nil {
		return err
	}

	// Create Kafka topics with specified configurations ...
	var topicConfigs []kafka.TopicConfig
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     partition,
			ReplicationFactor: kafkaClient.ReplicationFactor,
		})
	}

	if err := conn.CreateTopics(topicConfigs...); err != nil {
		return err
	}

	return nil
}

// setupHandler initializes the message handler for processing notifications.
func (s *Server) setupHandler() {
	// Initialize the handler with the logger, configuration, tracer, and Kafka producers ...
	s.handler = handler.NewHandler(s.logger, s.cfg, s.tracer.Tracer, s.cfg.ServiceClient.EmailService.Webhook, s.producerTopicMap, s.producerMap)
}

// createMessageProcessor creates a new instance of the message processor.
func (s *Server) createMessageProcessor() *mp.MessageProcessor {
	// Create a message processor using the provided components ...
	return mp.NewMessageProcessor(s.logger, s.cfg, s.handler, s.tracer.Tracer, s.producerTopicMap)
}

// prepareConsumerTopics prepares a list of consumer topics based on channel and priority.
func (s *Server) prepareConsumerTopics(channel string, priority string) []string {
	// Based on the provided channel and priority, prepare a list of consumer topics ...
	consumerTopics := strings.Split(s.cfg.KafkaTopic.Consumer, ",")
	switch priority {
	case "high":
		return s.createTopicList(consumerTopics, channel, constants.PRIORITY_HIGH)
	case "normal":
		return s.createTopicList(consumerTopics, channel, constants.PRIORITY_NORMAL)
	}
	return consumerTopics
}

// startConsumers starts Kafka consumers for the provided topics.
func (s *Server) startConsumers(ctx context.Context, consumerTopics []string, messageProcessor *mp.MessageProcessor) error {
	// Start Kafka consumers for each specified topic ...
	consumerBrokers := strings.Split(s.cfg.Kafka.ConsumerBroker, ",")
	fmt.Println("Consumers broker:", consumerBrokers)

	consumer := kafkaClient.NewConsumer(consumerBrokers)
	poolSize, err := strconv.Atoi(s.cfg.Kafka.PoolSize)
	if err != nil {
		return err
	}

	for _, topic := range consumerTopics {
		go consumer.StartWorkers(ctx, s.cfg.Kafka.GroupID, topic, poolSize, messageProcessor.ProcessMessage)
	}

	return nil
}

// createProducerMap creates a map of Kafka producers for different notification types.
func (s *Server) createProducerMap(producerTopics []string, producerBrokers []string) map[string]*kafkaClient.Producer {
	// Create Kafka producers for each specified notification type and associate them with topics ...
	producerMap := make(map[string]*kafkaClient.Producer)

	for _, topic := range producerTopics {
		producer := kafkaClient.NewProducer(producerBrokers, topic)

		// Associate each producer with its corresponding notification type ...
		if strings.Contains(topic, constants.NOTIF_TYPE_EMAIL) {
			producerMap[constants.NOTIF_TYPE_EMAIL] = producer
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_SMS_POOL) {
			producerMap[constants.NOTIF_TYPE_SMS_POOL] = producer
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_SMS) {
			producerMap[constants.NOTIF_TYPE_SMS] = producer
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_INAPP) {
			producerMap[constants.NOTIF_TYPE_INAPP] = producer
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_PUSH) {
			producerMap[constants.NOTIF_TYPE_PUSH] = producer
			continue
		}
	}

	return producerMap
}

// createTopicMap creates a map of topic names for easier access and association.
func (s *Server) createTopicMap(topics []string) map[string]string {
	// Create a mapping of topic names for easier access and association ...
	topicMap := make(map[string]string)

	for _, topic := range topics {
		// Associate each topic with its corresponding notification type ...
		if strings.Contains(topic, constants.NOTIF_TYPE_EMAIL) {
			topicMap[constants.NOTIF_TYPE_EMAIL] = topic
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_SMS_POOL) {
			topicMap[constants.NOTIF_TYPE_SMS_POOL] = topic
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_SMS) {
			topicMap[constants.NOTIF_TYPE_SMS] = topic
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_INAPP) {
			topicMap[constants.NOTIF_TYPE_INAPP] = topic
			continue
		} else if strings.Contains(topic, constants.NOTIF_TYPE_PUSH) {
			topicMap[constants.NOTIF_TYPE_PUSH] = topic
			continue
		}
	}

	return topicMap
}

// createTopicList creates a list of topic names with placeholders replaced by actual values.
func (s *Server) createTopicList(topics []string, channel, priority string) []string {
	// Create a list of topic names with placeholders replaced by actual values ...
	newTopics := []string{}

	for _, topic := range topics {
		topic = strings.ReplaceAll(topic, "<channel>", channel)
		topic = strings.ReplaceAll(topic, "<priority>", priority)
		newTopics = append(newTopics, topic)
	}

	return newTopics
}
