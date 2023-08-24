package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
)

func (h *Handler) HandlePush(ctx context.Context, parentMsg *model.PublishedKafkaMsg, pushMsg *model.Push, traceID string) error {
	ctx, span := h.tracer.Start(ctx, "Handler.HandlePush")
	defer span.End()

	pushMsg.Status = "created"
	pushBytes, err := json.Marshal(pushMsg)
	if err != nil {
		return err
	}

	parentMsg.Data = pushBytes

	if err := h.publishMessageToKafka(ctx, parentMsg, traceID, constants.NOTIF_TYPE_PUSH, pushMsg); err != nil {
		return fmt.Errorf("publish message to kafka failed: %w", err)
	}

	if _, err := h.sendPushMessageToProvider(ctx, pushMsg, traceID); err != nil {
		return fmt.Errorf("send message to profider failed: %w", err)
	}

	return nil
}

func (h *Handler) sendPushMessageToProvider(ctx context.Context, pushMsg *model.Push, traceID string) (string, error) {
	ctx, span := h.tracer.Start(ctx, "Handler.sendPushMessageToProvider")
	defer span.End()

	var msgId string
	var err error

	switch h.cfg.Project.Priority {
	case "high":
		h.logRestMessage(traceID, h.cfg.ServiceClient.OneSignalPushService.Url, constants.METHOD_POST, pushMsg, nil, "Sending push notification to Onesignal")

		_, msgId, err = h.sc.OneSignalPush(ctx, h.cfg.ServiceClient.OneSignalPushService.JmoAppId, pushMsg, traceID)
		if err != nil {
			h.logRestMessage(traceID, h.cfg.ServiceClient.OneSignalPushService.Url, constants.METHOD_POST, pushMsg, nil, "Error to send push notification")
			return "", err
		}

	case "normal":
		h.logRestMessage(traceID, h.cfg.ServiceClient.FcmPushService.Url, constants.METHOD_POST, pushMsg, nil, "Sending push notification to FCM")

		_, msgId, err = h.sc.FcmPush(ctx, pushMsg, traceID)
		if err != nil {
			h.logRestMessage(traceID, h.cfg.ServiceClient.OneSignalPushService.Url, constants.METHOD_POST, pushMsg, nil, "Error to send push notification")
			return "", err
		}
	}

	return msgId, nil
}
