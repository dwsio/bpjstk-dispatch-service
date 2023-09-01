package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
)

func (u *Usecase) HandlePush(ctx context.Context, parentMsg *model.PublishedKafkaMsg, pushMsg *model.Push) error {
	ctx, span := u.tracer.Start(ctx, "Usecase.HandlePush")
	defer span.End()

	pushMsg.Status = "created"
	pushBytes, err := json.Marshal(pushMsg)
	if err != nil {
		return err
	}

	parentMsg.Data = pushBytes

	if err := u.publishMessageToKafka(ctx, parentMsg, constants.NOTIF_TYPE_PUSH, pushMsg); err != nil {
		return fmt.Errorf("publish message to kafka failed: %w", err)
	}

	if _, err := u.sendPushMessageToProvider(ctx, pushMsg); err != nil {
		return fmt.Errorf("send message to provider failed: %w", err)
	}

	return nil
}

func (u *Usecase) sendPushMessageToProvider(ctx context.Context, pushMsg *model.Push) (string, error) {
	ctx, span := u.tracer.Start(ctx, "Usecase.sendPushMessageToProvider")
	defer span.End()

	var msgId string
	var err error

	switch u.cfg.Project.Priority {
	case "high":
		u.logRestMessage(ctx, u.cfg.ProviderClient.OneSignalPushProvider.Url, constants.METHOD_POST, pushMsg, nil, "Sending push notification to Onesignal")

		_, msgId, err = u.sc.OneSignalPush(ctx, u.cfg.ProviderClient.OneSignalPushProvider.JmoAppId, pushMsg)
		if err != nil {
			u.logRestMessage(ctx, u.cfg.ProviderClient.OneSignalPushProvider.Url, constants.METHOD_POST, pushMsg, nil, "Error to send push notification")
			return "", err
		}

	case "normal":
		u.logRestMessage(ctx, u.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, pushMsg, nil, "Sending push notification to FCM")

		_, msgId, err = u.sc.FcmPush(ctx, pushMsg)
		if err != nil {
			u.logRestMessage(ctx, u.cfg.ProviderClient.OneSignalPushProvider.Url, constants.METHOD_POST, pushMsg, nil, "Error to send push notification")
			return "", err
		}
	}

	return msgId, nil
}
