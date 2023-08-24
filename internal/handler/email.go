package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/utils"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
)

func (h *Handler) HandleEmail(ctx context.Context, parentMsg *model.PublishedKafkaMsg, emailMsg *model.Email, traceID string) error {
	ctx, span := h.tracer.Start(ctx, "Handler.HandleEmail")
	defer span.End()

	emailMsg.Status = "created"
	emailBytes, err := json.Marshal(emailMsg)
	if err != nil {
		return err
	}

	parentMsg.Data = emailBytes

	if err := h.publishMessageToKafka(ctx, parentMsg, traceID, constants.NOTIF_TYPE_EMAIL, emailMsg); err != nil {
		return fmt.Errorf("publish message to kafka failed: %w", err)
	}

	if err := h.sendEmailMessageToProvider(ctx, emailMsg, traceID); err != nil {
		return fmt.Errorf("send message to profider failed: %w", err)
	}

	return nil
}

func (h *Handler) sendEmailMessageToProvider(ctx context.Context, emailMsg *model.Email, traceID string) error {
	ctx, span := h.tracer.Start(ctx, "Handler.sendEmailMessageToProvider")
	defer span.End()

	utils.InjectWebhook(emailMsg, h.emailWebhookUrl)

	if emailMsg.Attachment != nil {
		attachments := []string{}

		for _, url := range emailMsg.Attachment {
			base64Data, err := utils.DownloadFileAsBase64(url)
			if err != nil {
				h.logKafkaMessage(traceID, h.producerTopicMap[constants.NOTIF_TYPE_EMAIL], emailMsg, err, "Error to download attachment")
				continue
			}

			attachments = append(attachments, base64Data)
		}

		emailMsg.Attachment = attachments
	}

	msg, _, err := h.sc.SendEmail(ctx, "noreply", emailMsg, traceID)
	if err != nil {
		return err
	}

	if !strings.EqualFold(msg, "Sukses") {
		return fmt.Errorf("send email failed with message: %s", msg)
	}

	return nil
}
