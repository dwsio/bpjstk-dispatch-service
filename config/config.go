package config

import (
	"os"

	kafkaClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/kafka"
	loggerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	metricClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/metric"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
)

type Config struct {
	Project        *Project             `mapstructure:"PROJECT"`
	Logger         *loggerClient.Config `mapstructure:"LOGGER_CLIENT"`
	Kafka          *kafkaClient.Config  `mapstructure:"KAFKA_CLIENT"`
	KafkaTopic     *KafkaTopic          `mapstructure:"KAFKA_TOPIC"`
	Tracer         *tracerClient.Config `mapstructure:"TRACER_CLIENT"`
	Metric         *metricClient.Config `mapstructure:"METRIC_CLIENT"`
	ProviderClient *ProviderClient      `mapstructure:"SERVICE_CLIENT"`
}

type Project struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Version     string `mapstructure:"VERSION"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Priority    string `mapstructure:"PRIORITY"`
	ServerIP    string `mapstructure:"SERVER_IP"`
}

type KafkaTopic struct {
	Producer string `mapstructure:"PRODUCER"`
	Consumer string `mapstructure:"CONSUMER"`
}

type ProviderClient struct {
	EmailProvider         *EmailProvider         `mapstructure:"EMAIL_PROVIDER"`
	SmsProvider           *SmsProvider           `mapstructure:"SMS_PROVIDER"`
	FcmPushProvider       *FcmPushProvider       `mapstructure:"FCM_PUSH_PROVIDER"`
	OneSignalPushProvider *OneSignalPushProvider `mapstructure:"ONESIGNAL_PUSH_PROVIDER"`
}

type EmailProvider struct {
	Url     string `mapstructure:"URL"`
	From    string `mapstructure:"FROM"`
	Webhook string `mapstructure:"WEBHOOK"`
}

type SmsProvider struct {
	Url      string `mapstructure:"URL"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

type FcmPushProvider struct {
	Url    string `mapstructure:"URL"`
	ApiKey string `mapstructure:"API_KEY"`
}

type OneSignalPushProvider struct {
	Url       string `mapstructure:"URL"`
	ApiKey    string `mapstructure:"API_KEY"`
	JmoAppId  string `mapstructure:"JMO_APP_ID"`
	SippAppId string `mapstructure:"SIPP_APP_ID"`
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
			Encoding: getEnv("LOGGER_ENCODING", "console"),
		},
		ProviderClient: &ProviderClient{
			EmailProvider: &EmailProvider{
				Url:     getEnv("EMAIL_PROVIDER_URL", "http://172.28.108.181:2014/WSCom/services/Main?wsdl"),
				From:    getEnv("EMAIL_PROVIDER_FROM", "noreply@bpjsketenagakerjaan.go.id"),
				Webhook: getEnv("EMAIL_PROVIDER_WEBHOOK", "http://172.28.108.245:8080/api/v1/webhook/"),
			},
			SmsProvider: &SmsProvider{
				Url:      getEnv("SMS_PROVIDER_URL", "http://172.28.108.181:2014/SmsApps/services/Main?wsdl"),
				Username: getEnv("SMS_PROVIDER_USERNAME", "sso"),
				Password: getEnv("SMS_PROVIDER_PASSWORD", "sso123"),
			},
			FcmPushProvider: &FcmPushProvider{
				Url:    getEnv("FCM_PUSH_PROVIDER_URL", "https://fcm.googleapis.com/fcm/send"),
				ApiKey: getEnv("FCM_PUSH_PROVIDER_API_KEY", ""),
			},
			OneSignalPushProvider: &OneSignalPushProvider{
				Url:       getEnv("ONESIGNAL_PUSH_PROVIDER_URL", "https://onesignal.com/api/v1/notifications"),
				ApiKey:    getEnv("ONESIGNAL_PUSH_PROVIDER_API_KEY", "MWE4M2U4OGEtMmRlZi00ODI0LTkxNDYtYjFiZmIyZTAzYzJk"),
				JmoAppId:  getEnv("ONESIGNAL_PUSH_PROVIDER_JMO_APP_ID", "40b2bca3-fbc3-47b1-a518-df6093404d7f"),
				SippAppId: getEnv("ONESIGNAL_PUSH_PROVIDER_SIPP_APP_ID", ""),
			},
		},
		Kafka: &kafkaClient.Config{
			ProducerBrokers: getEnv("KAFKA_PRODUCER_BROKERS", "localhost:29092"),
			ConsumerBrokers: getEnv("KAFKA_CONSUMER_BROKERS", "localhost:29093"),
			GroupID:         getEnv("KAFKA_GROUP_ID", "cns_dispatch_consumer"),
			PoolSize:        getEnv("KAFKA_POOL_SIZE", "10"),
			Partition:       getEnv("KAFKA_PARTITION", "10"),
		},
		KafkaTopic: &KafkaTopic{
			Producer: getEnv("KAFKA_TOPIC_PRODUCER", "cns_trc_email,cns_trc_sms,cns_trc_inapp,cns_trc_push,cns_trc_sms_pool"),
			Consumer: getEnv("KAFKA_TOPIC_CONSUMER", "cns_dsp_<channel>_email_<priority>,cns_dsp_<channel>_sms_<priority>,cns_dsp_<channel>_inapp_<priority>,cns_dsp_<channel>_push_<priority>"),
		},
		Tracer: &tracerClient.Config{
			Endpoint: getEnv("TRACER_ENDPOINT", "http://localhost:14268/api/traces"),
			Prefix:   getEnv("TRACER_PREFIX", "cns_dispatch"),
		},
		Metric: &metricClient.Config{
			Port:      getEnv("METRIC_PORT", ":8085"),
			Path:      getEnv("METRIC_PATH", "/metrics"),
			Prefix:    getEnv("METRIC_PREFIX", "cns_dispatch"),
			MeterName: getEnv("METRIC_METER_NAME", "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service"),
		},
	}
}
