package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"github.com/google/uuid"
)

func InjectMetadataToNewContext(parentCtx context.Context, traceID string, topic string) context.Context {
	ctx := context.WithValue(parentCtx, "trace_id", traceID)
	ctx = context.WithValue(ctx, "transaction_id", uuid.New().String())
	ctx = context.WithValue(ctx, "topic", topic)
	ctx = context.WithValue(ctx, "start_time", time.Now().UnixMilli())

	return ctx
}

func ExtractMetadataFromContext(ctx context.Context) (string, string, string, int64) {
	traceID := ctx.Value("trace_id").(string)
	transactionID := ctx.Value("transaction_id").(string)
	topic := ctx.Value("topic").(string)
	startTime := ctx.Value("start_time").(int64)

	return traceID, transactionID, topic, startTime
}

func StructToString(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func StringToStruct(jsonString string, data interface{}) error {
	if err := json.Unmarshal([]byte(jsonString), data); err != nil {
		return err
	}

	return nil
}

const maxSize = 10 * 1024 * 1024

func DownloadFileAsBase64(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status: %s", response.Status)
	}

	if response.ContentLength > maxSize {
		return "", fmt.Errorf("file size exceeds the limit")
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxSize))
	if err != nil {
		return "", err
	}

	base64Data := base64.StdEncoding.EncodeToString(body)

	return base64Data, nil
}

func InjectWebhook(emailMsg *model.Email, url string) error {
	msgBytes, err := json.Marshal(emailMsg)
	if err != nil {
		return err
	}

	// TODO: check again
	hookHtml := fmt.Sprintf(`<img src="%vdata=%v" alt="" width="1" height="1">`, url, msgBytes)

	emailMsg.ContentHTML += hookHtml

	return nil
}
