package config

import (
	"os"

	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	loggerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
)

type Config struct {
	Project       *Project
	Logger        *loggerClient.Config
	Kafka         *kafkaClient.Config
	KafkaTopic    *KafkaTopic
	Tracer        *tracerClient.Config
	ServiceClient *ServiceClient
}

type Project struct {
	ServiceName string
	Version     string
	Environment string
	Priority    string
	ServerIP    string
}

type KafkaTopic struct {
	Producer string
	Consumer string
}

type ServiceClient struct {
	EmailService         *EmailService
	SmsService           *SmsService
	FcmPushService       *FcmPushService
	OneSignalPushService *OneSignalPushService
}

type EmailService struct {
	Url     string
	From    string
	Webhook string
}

type SmsService struct {
	Url      string
	Username string
	Password string
}

type FcmPushService struct {
	Url    string
	ApiKey string
}

type OneSignalPushService struct {
	Url       string
	ApiKey    string
	JmoAppId  string
	SippAppId string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfigFromOS() *Config {
	return &Config{
		Project: &Project{
			ServiceName: getEnv("PROJECT_SERVICE_NAME", "cns-dispatch"),
			Version:     getEnv("PROJECT_VERSION", "v1.0.0"),
			Environment: getEnv("PROJECT_ENVIRONMENT", "dev"),
		},
		Logger: &loggerClient.Config{
			Encoding:   getEnv("LOGGER_ENCODING", "json"),
			Level:      getEnv("LOGGER_LEVEL", "info"),
			OutputPath: getEnv("LOGGER_OUTPUT_PATH", "stdout"),
			ErrorPath:  getEnv("LOGGER_ERROR_PATH", "stderr"),
		},
		ServiceClient: &ServiceClient{
			EmailService: &EmailService{
				Url:     getEnv("EMAIL_SERVICE_URL", "http://172.28.108.181:2014/WSCom/services/Main?wsdl"),
				From:    getEnv("EMAIL_SERVICE_FROM", "noreply@bpjsketenagakerjaan.go.id"),
				Webhook: getEnv("EMAIL_SERVICE_WEBHOOK", "http://172.28.108.245:8080/api/v1/webhook/"),
			},
			SmsService: &SmsService{
				Url:      getEnv("SMS_SERVICE_URL", "http://172.28.108.181:2014/SmsApps/services/Main?wsdl"),
				Username: getEnv("SMS_SERVICE_USERNAME", "sso"),
				Password: getEnv("SMS_SERVICE_PASSWORD", "sso123"),
			},
			FcmPushService: &FcmPushService{
				Url:    getEnv("FCM_PUSH_SERVICE_URL", "https://fcm.googleapis.com/fcm/send"),
				ApiKey: getEnv("FCM_PUSH_SERVICE_API_KEY", ""),
			},
			OneSignalPushService: &OneSignalPushService{
				Url:       getEnv("ONESIGNAL_PUSH_SERVICE_URL", "https://onesignal.com/api/v1/notifications"),
				ApiKey:    getEnv("ONESIGNAL_PUSH_SERVICE_API_KEY", "MWE4M2U4OGEtMmRlZi00ODI0LTkxNDYtYjFiZmIyZTAzYzJk"),
				JmoAppId:  getEnv("ONESIGNAL_PUSH_SERVICE_JMO_APP_ID", "40b2bca3-fbc3-47b1-a518-df6093404d7f"),
				SippAppId: getEnv("ONESIGNAL_PUSH_SERVICE_SIPP_APP_ID", ""),
			},
		},
		Kafka: &kafkaClient.Config{
			ProducerBroker: getEnv("KAFKA_PRODUCER_BROKER", "localhost:9095"),
			ConsumerBroker: getEnv("KAFKA_CONSUMER_BROKER", "localhost:9094"),
			GroupID:        getEnv("KAFKA_GROUP_ID", "cns_dispatch_consumer"),
			PoolSize:       getEnv("KAFKA_POOL_SIZE", "10"),
			Partition:      getEnv("KAFKA_PARTITION", "10"),
		},
		KafkaTopic: &KafkaTopic{
			Producer: getEnv("KAFKA_TOPIC_PRODUCER", "cns_trc_email,cns_trc_sms,cns_trc_inapp,cns_trc_push,cns_trc_sms_pool"),
			Consumer: getEnv("KAFKA_TOPIC_CONSUMER", "cns_dsp_<channel>_email_<priority>,cns_dsp_<channel>_sms_<priority>,cns_dsp_<channel>_inapp_<priority>,cns_dsp_<channel>_push_<priority>"),
		},
		Tracer: &tracerClient.Config{
			Endpoint: getEnv("TRACER_ENDPOINT", "http://localhost:14268/api/traces"),
		},
	}
}
