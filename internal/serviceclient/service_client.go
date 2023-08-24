package serviceclient

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

type IServiceClient interface {
	SendEmail(ctx context.Context, replyCfg string, email *model.Email, traceID string) (string, string, error)
	SendSms(ctx context.Context, sms *model.Sms, traceID string) (string, string, error)
	FcmPush(ctx context.Context, pushMsg *model.Push, traceID string) (*FcmPushRes, string, error)
	OneSignalPush(ctx context.Context, appId string, pushMsg *model.Push, traceID string) (*onesignal.CreateNotificationSuccessResponse, string, error)
}

type ServiceClient struct {
	logger          *logger.AppLogger
	cfg             *config.Config
	tracer          trace.Tracer
	onesignalClient *onesignal.APIClient
}

func NewServiceClient(logger *logger.AppLogger, cfg *config.Config, tracer trace.Tracer) *ServiceClient {
	return &ServiceClient{
		logger:          logger,
		cfg:             cfg,
		tracer:          tracer,
		onesignalClient: onesignal.NewAPIClient(onesignal.NewConfiguration()),
	}
}

func (sc *ServiceClient) NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 15 * time.Second,
			}).Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: 15 * time.Second,
			IdleConnTimeout:     30 * time.Second,
		},
	}
}
