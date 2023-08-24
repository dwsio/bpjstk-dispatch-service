package serviceclient

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	"github.com/OneSignal/onesignal-go-api"
)

func (sc *ServiceClient) OneSignalPush(ctx context.Context, appId string, pushMsg *model.Push, traceID string) (*onesignal.CreateNotificationSuccessResponse, string, error) {
	ctx, span := sc.tracer.Start(ctx, "ServiceClient.OneSignalPush")
	defer span.End()

	notification := *onesignal.NewNotification(appId)
	notification.SetIncludePlayerIds(pushMsg.PlayerIds)
	notification.SetHeadings(onesignal.StringMap{En: &pushMsg.Heading})
	notification.SetContents(onesignal.StringMap{En: &pushMsg.Content})
	notification.SetBigPicture(pushMsg.PictureUrl)
	notification.SetIsIos(pushMsg.IsIos)

	appAuth := context.WithValue(ctx, onesignal.AppAuth, sc.cfg.ServiceClient.OneSignalPushService.ApiKey)

	notifSuccesRes, httpRes, err := sc.onesignalClient.DefaultApi.CreateNotification(appAuth).Notification(notification).Execute()
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.OneSignalPushService.Url, constants.METHOD_POST, pushMsg, nil, "Error to send push notification", constants.StatusAborted)
		return nil, "", err
	}

	httpResBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.OneSignalPushService.Url, constants.METHOD_POST, "", err, "Get result body", "Failed to get result body")
		return nil, "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.OneSignalPushService.Url, constants.METHOD_POST, string(httpResBody), fmt.Errorf("status code %v", httpRes.StatusCode), "Check HTTP result code", "Failed to send request")
		return nil, "", err
	}

	sc.logRestMessage(traceID, sc.cfg.ServiceClient.OneSignalPushService.Url, constants.METHOD_POST, string(httpResBody), nil, "Success to send sms request", "")

	return notifSuccesRes, notifSuccesRes.Id, nil
}
