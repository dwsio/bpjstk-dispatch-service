package provider_client

import (
	"context"
	"fmt"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	contextMD "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/context_metadata"
	loggerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
)

func (pc *ProviderClient) logRestMessage(ctx context.Context, endpoint, method string, data interface{}, err error, activity string) {
	logFields := pc.getRestLogFields(ctx, endpoint, method, data, err, activity)
	if err != nil {
		logFields.LogLevel = constants.LEVEL_ERROR
		logFields.Error = constants.TRUE
		logFields.Message.ErrorMessage = err.Error()
	}
	pc.logger.StructuredPrint(logFields)
}

func (pc *ProviderClient) getRestLogFields(ctx context.Context, endpoint, method string, data interface{}, err error, activity string) *loggerClient.LogFields {
	metadata, _ := contextMD.GetMetadataFromContext(ctx)

	return &loggerClient.LogFields{
		Timestamp:     time.Now().Format(constants.TIME_LAYOUT_FORMAT),
		LogLevel:      getLoggerLogLevel(err),
		TransactionID: fmt.Sprintf("TR%s", metadata.TransactionID),
		ServiceName:   pc.cfg.Project.ServiceName,
		Endpoint:      endpoint,
		Protocol:      constants.PROTOCOL_HTTP,
		MethodType:    method,
		ExecutionType: constants.ASYNC,
		ContentType:   constants.CONTENT_JSON,
		FunctionName:  "",
		UserInfo: &loggerClient.UserInfo{
			Username: "",
			Role:     "",
			Others:   "",
		},
		ExecutionTime:     fmt.Sprint(time.Since(time.UnixMilli(metadata.StartTime)).Milliseconds()),
		ServerIP:          pc.cfg.Project.ServerIP,
		ClientIP:          "",
		EventName:         "",
		TraceID:           fmt.Sprintf("TC%s", metadata.TraceID),
		PrevTransactionID: "",
		Body:              fmt.Sprint(data),
		Result:            "",
		Error:             constants.FALSE,
		FlagStartOrStop:   constants.STOP,
		Message: &loggerClient.Message{
			Activity:          activity,
			ObjectPerformedOn: "",
			ResultOfActivity:  getLoggerResult(err),
			ErrorCode:         "",
			ErrorMessage:      "",
			ShortDescription:  "",
		},
	}
}

func getLoggerLogLevel(err error) string {
	if err != nil {
		return constants.LEVEL_ERROR
	}
	return constants.LEVEL_INFO
}

func getLoggerResult(err error) string {
	if err != nil {
		return constants.StatusAborted
	}
	return ""
}
