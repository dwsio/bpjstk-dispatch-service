package provider_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
)

const (
	MAX_TTL         = 2419200
	Priority_HIGH   = "high"
	Priority_NORMAL = "normal"
)

type FcmPushReq struct {
	Data                  map[string]interface{} `json:"data,omitempty"`
	To                    string                 `json:"to,omitempty"`
	RegistrationIds       []string               `json:"registration_ids,omitempty"`
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Priority              string                 `json:"priority,omitempty"`
	Notification          NotificationPayload    `json:"notification,omitempty"`
	ContentAvailable      bool                   `json:"content_available,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
	Condition             string                 `json:"condition,omitempty"`
	MutableContent        bool                   `json:"mutable_content,omitempty"`
}

type NotificationPayload struct {
	Title            string `json:"title,omitempty"`
	Body             string `json:"body,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Image            string `json:"image,omitempty"`
	Sound            string `json:"sound,omitempty"`
	Badge            string `json:"badge,omitempty"`
	Tag              string `json:"tag,omitempty"`
	Color            string `json:"color,omitempty"`
	ClickAction      string `json:"click_action,omitempty"`
	BodyLocKey       string `json:"body_loc_key,omitempty"`
	BodyLocArgs      string `json:"body_loc_args,omitempty"`
	TitleLocKey      string `json:"title_loc_key,omitempty"`
	TitleLocArgs     string `json:"title_loc_args,omitempty"`
	AndroidChannelID string `json:"android_channel_id,omitempty"`
}

type FcmPushRes struct {
	Ok           bool
	StatusCode   int
	MulticastId  int64               `json:"multicast_id"`
	Success      int                 `json:"success"`
	Fail         int                 `json:"failure"`
	CanonicalIds int                 `json:"canonical_ids"`
	Results      []map[string]string `json:"results,omitempty"`
	MessageId    int64               `json:"message_id,omitempty"`
	Error        string              `json:"error,omitempty"`
	RetryAfter   string
}

func (pc *ProviderClient) FcmPush(ctx context.Context, pushMsg *model.Push) (*FcmPushRes, string, error) {
	_, span := pc.tracer.Start(ctx, "ProviderClient.FCMPush")
	defer span.End()

	client := pc.NewHttpClient()

	httpReqBody, err := buildPushReqBody(pushMsg)
	if err != nil {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, pushMsg, err, "Create request body")
		return nil, "", err
	}

	httpReq, err := http.NewRequest(constants.METHOD_POST, pc.cfg.ProviderClient.FcmPushProvider.Url, bytes.NewBuffer(httpReqBody))
	if err != nil {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, string(httpReqBody), err, "Create new request")
		return nil, "", err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("key=%v", pc.cfg.ProviderClient.FcmPushProvider.ApiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	httpRes, err := client.Do(httpReq)
	if err != nil {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, string(httpReqBody), err, "Get http result")
		return nil, "", err
	}
	defer httpRes.Body.Close()

	httpResBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, string(httpReqBody), err, "Get result body")
		return nil, "", err
	}

	fcmPushRes := new(FcmPushRes)
	fcmPushRes.StatusCode = httpRes.StatusCode

	if fcmPushRes.StatusCode != http.StatusOK {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, string(httpResBody), fmt.Errorf("status code %v", fcmPushRes.StatusCode), "Check HTTP result code")
		return nil, "", err
	}

	if err := json.Unmarshal(httpResBody, &fcmPushRes); err != nil {
		pc.logRestMessage(ctx, pc.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, string(httpResBody), err, "Unmarshal response body")
		return nil, "", err
	}

	pc.logRestMessage(ctx, pc.cfg.ProviderClient.FcmPushProvider.Url, constants.METHOD_POST, string(httpResBody), nil, "Success to send sms request")

	return fcmPushRes, fmt.Sprint(fcmPushRes.MessageId), nil
}

func buildPushReqBody(pushMsg *model.Push) ([]byte, error) {
	body := &FcmPushReq{
		Data:             pushMsg.Data,
		RegistrationIds:  pushMsg.PlayerIds,
		Priority:         Priority_NORMAL,
		ContentAvailable: true,
		TimeToLive:       MAX_TTL,
		Notification: NotificationPayload{
			Title: pushMsg.Heading,
			Body:  pushMsg.Content,
			Image: pushMsg.PictureUrl,
		},
	}

	return json.Marshal(body)
}
