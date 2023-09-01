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

func (u *Usecase) HandleSMS(ctx context.Context, parentMsg *model.PublishedKafkaMsg, smsMsg *model.Sms) error {
	ctx, span := u.tracer.Start(ctx, "Usecase.HandleSMS")
	defer span.End()

	smsMsg.Status = "created"
	smsBytes, err := json.Marshal(smsMsg)
	if err != nil {
		return err
	}

	parentMsg.Data = smsBytes

	if err := u.publishMessageToKafka(ctx, parentMsg, constants.NOTIF_TYPE_SMS, smsMsg); err != nil {
		return fmt.Errorf("publish message to kafka failed: %w", err)
	}

	kodeStruct, err := u.sendSMSMessageToProvider(ctx, smsMsg)
	if err != nil {
		return fmt.Errorf("send message to provider failed: %w", err)
	}

	smsMsg.Status = "on process"
	smsMsg.MessageId = kodeStruct.MsgID

	smsBytes, err = json.Marshal(smsMsg)
	if err != nil {
		return err
	}

	parentMsg.Data = smsBytes

	if err := u.publishMessageToKafka(ctx, parentMsg, constants.NOTIF_TYPE_SMS_POOL, smsMsg); err != nil {
		return fmt.Errorf("publish message to kafka failed: %w", err)
	}

	return nil
}

func (u *Usecase) sendSMSMessageToProvider(ctx context.Context, smsMsg *model.Sms) (*model.SendSmsResCode, error) {
	ctx, span := u.tracer.Start(ctx, "Usecase.sendSMSMessageToProvider")
	defer span.End()

	msg, kode, err := u.sc.SendSms(ctx, smsMsg)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(msg, "SUCCESS") {
		return nil, fmt.Errorf("send sms failed with message: %s", msg)
	}

	kodeStruct := model.SendSmsResCode{}
	if err := utils.StringToStruct(kode, &kodeStruct); err != nil {
		return nil, err
	}

	return &kodeStruct, nil
}
