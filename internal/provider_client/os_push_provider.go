package provider_client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	"github.com/OneSignal/onesignal-go-api"
)

func (pc *ProviderClient) OneSignalPush(ctx context.Context, appId string, pushMsg *model.Push) (*onesignal.CreateNotificationSuccessResponse, string, error) {
	ctx, span := pc.tracer.Start(ctx, "ProviderClient.OneSignalPush")
	defer span.End()

	notification := *onesignal.NewNotification(appId)
	notification.SetIncludePlayerIds(pushMsg.PlayerIds)
	notification.SetHeadings(onesignal.StringMap{En: &pushMsg.Heading})
	notification.SetContents(onesignal.StringMap{En: &pushMsg.Content})
	notification.SetBigPicture(pushMsg.PictureUrl)
	notification.SetIsIos(pushMsg.IsIos)

	appAuth := context.WithValue(ctx, onesignal.AppAuth, pc.cfg.ProviderClient.OneSignalPushProvider.ApiKey)

	notifSuccesRes, httpRes, err := pc.onesignalClient.DefaultApi.CreateNotification(appAuth).Notification(notification).Execute()
	if err != nil {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.OneSignalPushProvider.Url, constants.METHOD_POST, pushMsg, nil, "Error to send push notification")
		return nil, "", err
	}

	httpResBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.OneSignalPushProvider.Url, constants.METHOD_POST, "", err, "Get result body")
		return nil, "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.OneSignalPushProvider.Url, constants.METHOD_POST, string(httpResBody), fmt.Errorf("status code %v", httpRes.StatusCode), "Check HTTP result code")
		return nil, "", err
	}

	pc.logRestMessage(ctx, pc.cfg.ProviderClient.OneSignalPushProvider.Url, constants.METHOD_POST, string(httpResBody), nil, "Success to send sms request")

	return notifSuccesRes, notifSuccesRes.Id, nil
}
