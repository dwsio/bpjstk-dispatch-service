package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.LoadConfigFromOS()
	log := logger.NewAppLogger()

	at, err := tracerClient.NewAppTracer(ctx, cfg.Tracer, cfg.Project.ServiceName, cfg.Project.Version)
	if err != nil {
		panic(err)
	}

	defer func() error {
		if err := at.TraceProvider.Shutdown(ctx); err != nil {
			panic(err)
		}

		return nil
	}()

	tracer := at.Tracer

	producerBroker := strings.Split(cfg.Kafka.ConsumerBrokers, ",")
	topics := strings.Split("cns_dsp_jmo_email_reg,cns_dsp_jmo_sms_reg,cns_dsp_jmo_inapp_reg,cns_dsp_jmo_push_reg", ",")
	fmt.Printf("Brokers: %v Topics: %+v\n", producerBroker, topics)

	kafkaConn, err := kafka.DialContext(ctx, "tcp", producerBroker[0])
	if err != nil {
		panic(err)
	}
	defer kafkaConn.Close()

	initKafkaTopics(ctx, kafkaConn, topics)

	var producers []*kafkaClient.Producer

	for _, v := range topics {
		producer := kafkaClient.NewProducer(producerBroker, v)
		defer producer.Close()

		producers = append(producers, producer)
	}

	wg := sync.WaitGroup{}

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	emailToKafka(ctx, cfg, log, tracer, producers[0])
	// }()

	wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	smsToKafka(ctx, cfg, log, tracer, producers[1])
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	pushToKafka(ctx, cfg, log, tracer, producers[3])
	// }()

	wg.Wait()

	<-ctx.Done()
}

// func emailToKafka(ctx context.Context, cfg *config.Config, log *logger.AppLogger, tracer trace.Tracer, kafkaProducer *kafkaClient.Producer) {
// 	startTime := time.Now()

// 	for i := 1; i < 10; i++ {
// 		newCtx, span := tracer.Start(context.Background(), "test.emailToKafka")
// 		defer span.End()

// 		emailMsg := &model.Email{
// 			RecipientTo:  []string{"dw001423@abishar.com"},
// 			RecipientCc:  []string{"dw001423@abishar.com"},
// 			RecipientBcc: []string{},
// 			Subject:      "Test",
// 			ContentText:  "Test",
// 			IsHTML:       true,
// 			ContentHTML:  "<p>Test</p>",
// 			IsAttach:     false,
// 		}
// 		emailBytes, _ := json.Marshal(emailMsg)

// 		fromPreferenceKafkaMsg := &model.ConsumedKafkaMsg{
// 			CategoryName: "Email",
// 			TypeName:     "Campaign",
// 			ChannelName:  "jmo",
// 			Data:         emailBytes,
// 		}
// 		msgBytes, _ := json.Marshal(fromPreferenceKafkaMsg)

// 		kafkaMsg := kafka.Message{
// 			Value: msgBytes,
// 			Time:  time.Now().UTC(),
// 		}

// 		tracerClient.InjectKafkaTracingHeadersToCarrier(newCtx, tracerClient.NewKafkaHeadersCarrier(&kafkaMsg.Headers))

// 		injectTraceIDToKafkaHeaders(&kafkaMsg.Headers, uuid.New().String())
// 		fmt.Println("Sending email message...")

// 		time.Sleep(10 * time.Second)

// 		if err := kafkaProducer.PublishMessage(newCtx, tracer, kafkaMsg); err != nil {
// 			fmt.Printf("Publish send email message FAILED %v\n", err)
// 		}
// 	}
// 	fmt.Printf("time elapsed: %v", time.Since(startTime).Seconds())
// }

func smsToKafka(ctx context.Context, cfg *config.Config, log *logger.AppLogger, tracer trace.Tracer, kafkaProducer *kafkaClient.Producer) {
	startTime := time.Now()

	for i := 1; i < 15; i++ {
		newCtx, span := tracer.Start(context.Background(), "test.smsToKafka")
		defer span.End()

		smsMsg := &model.Sms{
			RecipientPhoneNumber: "+6281212763660",
			Content:              "Test CNS",
		}
		smsBytes, _ := json.Marshal(smsMsg)

		parentMsg := &model.ConsumedKafkaMsg{
			CategoryName: "sms",
			TypeName:     "Otp",
			ChannelName:  "jmo",
			Data:         smsBytes,
		}
		msgBytes, _ := json.Marshal(parentMsg)

		kafkaMsg := kafka.Message{
			Value: msgBytes,
			Time:  time.Now().UTC(),
		}

		tracerClient.InjectKafkaTracingHeadersToCarrier(newCtx, &kafkaMsg.Headers)

		fmt.Println("context", newCtx)

		traceID := uuid.New().String()
		injectTraceIDToKafkaHeaders(&kafkaMsg.Headers, traceID)

		fmt.Println("trace id", traceID)

		time.Sleep(10 * time.Second)

		if err := kafkaProducer.PublishMessage(newCtx, kafkaMsg); err != nil {
			fmt.Printf("Publish send SMS message FAILED %v\n", err)
		}
	}

	fmt.Printf("time elapsed: %v\n", time.Since(startTime).Seconds())
}

// func pushToKafka(ctx context.Context, cfg *config.Config, log *logger.AppLogger, tracer trace.Tracer, kafkaProducer *kafkaClient.Producer) {

// 	startTime := time.Now()

// 	for i := 1; i < 10; i++ {
// 		newCtx, span := tracer.Start(context.Background(), "test.pushToKafka")
// 		defer span.End()

// 		pushMsg := &model.Push{
// 			PlayerIds:  []string{"18b2fdb8-be46-484a-81d1-7c8a19e6e9f7"},
// 			Heading:    "Heading",
// 			Content:    "Content",
// 			PictureUrl: "https://sendbird.imgix.net/cms/20220328_Push-Notification-Testing-Tool_Blog-Cover.png",
// 			IsIos:      false,
// 		}
// 		pushBytes, _ := json.Marshal(pushMsg)

// 		fromPreferenceKafkaMsg := &model.ConsumedKafkaMsg{
// 			CategoryName: "push",
// 			TypeName:     "Campaign",
// 			ChannelName:  "jmo",
// 			Data:         pushBytes,
// 		}

// 		msgBytes, _ := json.Marshal(fromPreferenceKafkaMsg)

// 		kafkaMsg := kafka.Message{
// 			Value: msgBytes,
// 			Time:  time.Now().UTC(),
// 		}
// tracerClient.InjectKafkaTracingHeadersToCarrier(newCtx, tracerClient.NewKafkaHeadersCarrier(&kafkaMsg.Headers))

// 		injectTraceIDToKafkaHeaders(&kafkaMsg.Headers, uuid.New().String())

// 		fmt.Println("Sending push message...")

// 		time.Sleep(10 * time.Second)

// 		if err := kafkaProducer.PublishMessage(newCtx, tracer, kafkaMsg); err != nil {
// 			fmt.Printf("Publish send push message FAILED %v\n", err)
// 		}
// 	}
// 	fmt.Printf("time elapsed: %v", time.Since(startTime).Seconds())
// }

func initKafkaTopics(ctx context.Context, kafkaConn *kafka.Conn, topics []string) error {
	var topicConfigs []kafka.TopicConfig
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     10,
			ReplicationFactor: 1,
		})
	}

	if err := kafkaConn.CreateTopics(topicConfigs...); err != nil {
		return errors.Wrap(err, "conn.CreateTopics")
	}

	return nil
}

func injectTraceIDToKafkaHeaders(headers *[]kafka.Header, traceID string) {
	*headers = append(*headers, kafka.Header{
		Key:   "trace_id",
		Value: []byte(traceID),
	})
}
