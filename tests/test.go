package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

type contextKey string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	traceIDKey := contextKey("trace_id")

	ctx = context.WithValue(ctx, traceIDKey, "test123")

	cfg := config.LoadConfigFromOS()

	log, _ := logger.NewAppLogger(cfg.Logger)

	/**
	 *  OpenTelemetry Tracer
	 **/

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

	ctx, span := tracer.Start(ctx, "test.main")
	defer span.End()

	producerBroker := strings.Split(cfg.Kafka.ConsumerBroker, ",")
	topics := strings.Split("cns_dsp_jmo_email_reg,cns_dsp_jmo_sms_reg,cns_dsp_jmo_inapp_reg,cns_dsp_jmo_push_reg", ",")
	fmt.Printf("Brokers: %v Topics: %+v\n", producerBroker, topics)

	kafkaConn, err := connectKafkaBrokers(ctx, producerBroker[0])
	if err != nil {
		panic(err)
	}
	defer kafkaConn.Close() // nolint: errcheck

	initKafkaTopics(ctx, kafkaConn, topics)

	var producers []*kafkaClient.Producer

	for _, v := range topics {
		producer := kafkaClient.NewProducer(producerBroker, v)
		defer producer.Close()

		producers = append(producers, producer)
	}

	// emailToKafka(ctx, cfg, log, producers[0], span)
	smsToKafka(ctx, cfg, log, tracer, producers[1])
	// pushToKafka(ctx, cfg, log, producers[3], span)
}

// func emailToKafka(ctx context.Context, cfg *config.Config, log logger.ILogger, kafkaProducer *kafkaClient.Producer, span opentracing.Span) {
// 	emailMsg := &model.Email{
// 		RecipientTo:  []string{"dw001423@abishar.com"},
// 		RecipientCc:  []string{"dw001423@abishar.com"},
// 		RecipientBcc: []string{},
// 		Subject:      "Test",
// 		Payload: &model.EmailPayload{
// 			ContentText: "Test",
// 			IsHTML:      true,
// 			ContentHTML: "<p>Test</p>",
// 			IsAttach:    false,
// 		},
// 	}
// 	emailBytes, _ := json.Marshal(emailMsg)

// 	fromPreferenceKafkaMsg := &model.ConsumedKafkaMsg{
// 		CategoryName: "Email",
// 		TypeName:     "Campaign",
// 		ChannelName:  "jmo",
// 		Data:         emailBytes,
// 	}
// 	msgBytes, _ := json.Marshal(fromPreferenceKafkaMsg)

// 	kafkaMsg := kafka.Message{
// 		Value:   msgBytes,
// 		Time:    time.Now().UTC(),
// 		Headers: jaeger.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
// 	}

// 	kafkaMsg.Headers = append(kafkaMsg.Headers, kafka.Header{
// 		Key:   "trace_id",
// 		Value: []byte("test123"),
// 	})

// 	startTime := time.Now()
// 	wg := &sync.WaitGroup{}

// 	for i := 1; i < 2; i++ {
// 		wg.Add(1)

// 		go func(i int) {
// 			defer wg.Done()
// 			printSend(i)

// 			if err := kafkaProducer.PublishMessage(ctx, kafkaMsg); err != nil {
// 				fmt.Printf("Publish send email message FAILED %v\n", err)
// 			}
// 		}(i)

// 	}
// 	wg.Wait()
// 	fmt.Printf("time elapsed: %v", time.Since(startTime).Seconds())
// }

func smsToKafka(ctx context.Context, cfg *config.Config, log *logger.AppLogger, tracer trace.Tracer, kafkaProducer *kafkaClient.Producer) {
	ctx, span := tracer.Start(ctx, "test.smsToKafka")
	defer span.End()

	smsMsg := &model.Sms{
		RecipientPhoneNumber: "+6281212763660",
		Content:              "Test CNS",
	}
	smsBytes, _ := json.Marshal(smsMsg)

	fromPreferenceKafkaMsg := &model.ConsumedKafkaMsg{
		CategoryName: "sms",
		TypeName:     "Otp1",
		ChannelName:  "jmo",
		Data:         smsBytes,
	}
	msgBytes, _ := json.Marshal(fromPreferenceKafkaMsg)

	kafkaMsg := kafka.Message{
		Value: msgBytes,
		Time:  time.Now().UTC(),
	}

	carrier := tracerClient.NewKafkaHeadersCarrier(&kafkaMsg.Headers)

	tracerClient.InjectKafkaTracingHeadersToCarrier(ctx, carrier)

	kafkaMsg.Headers = append(kafkaMsg.Headers, kafka.Header{
		Key:   "trace_id",
		Value: []byte("test123"),
	})

	startTime := time.Now()
	wg := &sync.WaitGroup{}

	for i := 1; i < 2; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			printSend(i)

			if err := kafkaProducer.PublishMessage(ctx, tracer, kafkaMsg); err != nil {
				fmt.Printf("Send message failed %v\n", err)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("time elapsed: %v\n", time.Since(startTime).Seconds())
}

// func pushToKafka(ctx context.Context, cfg *config.Config, log logger.ILogger, kafkaProducer *kafkaClient.Producer, span opentracing.Span) {
// 	pushMsg := &model.Push{
// 		PlayerIds:  []string{"18b2fdb8-be46-484a-81d1-7c8a19e6e9f7"},
// 		Heading:    "Heading",
// 		Content:    "Content",
// 		PictureUrl: "https://sendbird.imgix.net/cms/20220328_Push-Notification-Testing-Tool_Blog-Cover.png",
// 		IsIos:      false,
// 	}
// 	pushBytes, _ := json.Marshal(pushMsg)

// 	fromPreferenceKafkaMsg := &model.ConsumedKafkaMsg{
// 		CategoryName: "push",
// 		TypeName:     "Campaign",
// 		ChannelName:  "jmo",
// 		Data:         pushBytes,
// 	}

// 	msgBytes, _ := json.Marshal(fromPreferenceKafkaMsg)

// 	kafkaMsg := kafka.Message{
// 		Value:   msgBytes,
// 		Time:    time.Now().UTC(),
// 		Headers: jaeger.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
// 	}

// 	kafkaMsg.Headers = append(kafkaMsg.Headers, kafka.Header{
// 		Key:   "trace_id",
// 		Value: []byte("test123"),
// 	})

// 	startTime := time.Now()
// 	wg := &sync.WaitGroup{}

// 	for i := 1; i < 11; i++ {
// 		wg.Add(1)

// 		go func(i int) {
// 			defer wg.Done()
// 			printSend(i)

// 			if err := kafkaProducer.PublishMessage(ctx, kafkaMsg); err != nil {
// 				fmt.Printf("Publish send push message FAILED %v\n", err)
// 			}
// 		}(i)
// 	}
// 	wg.Wait()
// 	fmt.Printf("time elapsed: %v\n", time.Since(startTime).Seconds())
// }

func printSend(i int) {
	fmt.Printf("Sending message...%v\n", i)
}

func connectKafkaBrokers(ctx context.Context, broker string) (*kafka.Conn, error) {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, broker)
	if err != nil {
		return nil, errors.Wrap(err, "kafka.NewKafkaCon")
	}

	_, err = kafkaConn.Brokers()
	if err != nil {
		return nil, errors.Wrap(err, "kafkaConn.Brokers")
	}

	return kafkaConn, nil
}

func initKafkaTopics(ctx context.Context, kafkaConn *kafka.Conn, topics []string) error {
	controller, err := kafkaConn.Controller()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Controller")
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		return errors.Wrap(err, "kafka.DialContext")
	}
	defer conn.Close() // nolint: errcheck

	var topicConfigs []kafka.TopicConfig
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     10,
			ReplicationFactor: 1,
		})
	}

	if err := conn.CreateTopics(topicConfigs...); err != nil {
		return errors.Wrap(err, "conn.CreateTopics")
	}

	return nil
}
