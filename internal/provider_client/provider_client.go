package provider_client

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"

	"github.com/OneSignal/onesignal-go-api"
	"go.opentelemetry.io/otel/trace"
)

type IProviderClient interface {
	SendEmail(ctx context.Context, replyCfg string, email *model.Email) (string, string, error)
	SendSms(ctx context.Context, sms *model.Sms) (string, string, error)
	FcmPush(ctx context.Context, pushMsg *model.Push) (*FcmPushRes, string, error)
	OneSignalPush(ctx context.Context, appId string, pushMsg *model.Push) (*onesignal.CreateNotificationSuccessResponse, string, error)
}

type ProviderClient struct {
	logger          *logger.AppLogger
	cfg             *config.Config
	tracer          trace.Tracer
	onesignalClient *onesignal.APIClient
}

func NewProviderClient(logger *logger.AppLogger, cfg *config.Config, tracer trace.Tracer) *ProviderClient {
	return &ProviderClient{
		logger:          logger,
		cfg:             cfg,
		tracer:          tracer,
		onesignalClient: onesignal.NewAPIClient(onesignal.NewConfiguration()),
	}
}

func (pc *ProviderClient) NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: 10 * time.Second,
			IdleConnTimeout:     30 * time.Second,
		},
	}
}
