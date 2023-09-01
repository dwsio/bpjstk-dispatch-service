package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/utils"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
)

func (u *Usecase) HandleEmail(ctx context.Context, parentMsg *model.PublishedKafkaMsg, emailMsg *model.Email) error {
	ctx, span := u.tracer.Start(ctx, "Usecase.HandleEmail")
	defer span.End()

	emailMsg.Status = "created"
	emailBytes, err := json.Marshal(emailMsg)
	if err != nil {
		return err
	}

	parentMsg.Data = emailBytes

	if err := u.publishMessageToKafka(ctx, parentMsg, constants.NOTIF_TYPE_EMAIL, emailMsg); err != nil {
		return fmt.Errorf("publish message to kafka failed: %w", err)
	}

	if err := u.sendEmailMessageToProvider(ctx, emailMsg); err != nil {
		return fmt.Errorf("send message to provider failed: %w", err)
	}

	return nil
}

func (u *Usecase) sendEmailMessageToProvider(ctx context.Context, emailMsg *model.Email) error {
	ctx, span := u.tracer.Start(ctx, "Usecase.sendEmailMessageToProvider")
	defer span.End()

	utils.InjectWebhook(emailMsg, u.emailWebhookUrl)

	if emailMsg.Attachment != nil {
		attachments := []string{}

		for _, url := range emailMsg.Attachment {
			base64Data, err := utils.DownloadFileAsBase64(url)
			if err != nil {
				u.logKafkaMessage(ctx, emailMsg, err, "Error to download attachment")
				continue
			}

			attachments = append(attachments, base64Data)
		}

		emailMsg.Attachment = attachments
	}

	msg, _, err := u.sc.SendEmail(ctx, "noreply", emailMsg)
	if err != nil {
		return err
	}

	if !strings.EqualFold(msg, "Sukses") {
		return fmt.Errorf("send email failed with message: %s", msg)
	}

	return nil
}
