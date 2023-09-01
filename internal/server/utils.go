package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	messageProcessor "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/message_processor"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	"github.com/segmentio/kafka-go"
)

func (s *Server) initKafkaTopics(ctx context.Context, producerBroker string, topics []string) error {
	conn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return err
	}

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

func (s *Server) prepareConsumerTopics(channel string, priority string) []string {
	consumerTopics := strings.Split(s.cfg.KafkaTopic.Consumer, ",")
	switch priority {
	case "high":
		return s.createTopicList(consumerTopics, channel, constants.PRIORITY_HIGH)
	case "normal":
		return s.createTopicList(consumerTopics, channel, constants.PRIORITY_NORMAL)
	}
	return consumerTopics
}

func (s *Server) startConsumers(ctx context.Context, consumerTopics []string, messageProcessor *messageProcessor.MessageProcessor) error {
	consumerBrokers := strings.Split(s.cfg.Kafka.ConsumerBrokers, ",")
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

func (s *Server) createProducerMap(producerTopics []string, producerBrokers []string) map[string]*kafkaClient.Producer {
	producerMap := make(map[string]*kafkaClient.Producer)

	for _, topic := range producerTopics {
		producer := kafkaClient.NewProducer(producerBrokers, topic)

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

func (s *Server) createTopicMap(topics []string) map[string]string {
	topicMap := make(map[string]string)

	for _, topic := range topics {
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

func (s *Server) createTopicList(topics []string, channel, priority string) []string {

	newTopics := []string{}

	for _, topic := range topics {
		topic = strings.ReplaceAll(topic, "<channel>", channel)
		topic = strings.ReplaceAll(topic, "<priority>", priority)
		newTopics = append(newTopics, topic)
	}

	return newTopics
}
