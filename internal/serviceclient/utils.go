package serviceclient

import (
	"fmt"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
)

func (sc *ServiceClient) logRestMessage(traceID string, endpoint string, method string, reqBody interface{}, err error, act string, res string) {
	logFields := &logger.LogFields{
		Timestamp:     time.Now().Format(constants.TIME_LAYOUT_FORMAT),
		LogLevel:      constants.LEVEL_INFO,
		TransactionID: "",
		ServiceName:   sc.cfg.Project.ServiceName,
		Endpoint:      endpoint,
		Protocol:      constants.PROTOCOL_HTTP,
		MethodType:    method,
		ExecutionType: constants.SYNC,
		ContentType:   constants.CONTENT_JSON,
		FunctionName:  "",
		UserInfo: &logger.UserInfo{
			Username: "",
			Role:     "",
			Others:   "",
		},
		ExecutionTime:     "",
		ServerIP:          sc.cfg.Project.ServerIP,
		ClientIP:          "",
		EventName:         "",
		TraceID:           traceID,
		PrevTransactionID: "",
		Body:              fmt.Sprintf("%v", reqBody),
		Result:            "",
		Error:             constants.FALSE,
		FlagStartOrStop:   constants.STOP,
		Message: &logger.Message{
			Activity:          act,
			ObjectPerformedOn: "",
			ResultOfActivity:  res,
			ErrorCode:         "",
			ErrorMessage:      "",
			ShortDescription:  "",
		},
	}

	if err != nil {
		logFields.LogLevel = constants.LEVEL_ERROR
		logFields.Error = constants.TRUE
		logFields.Message.ErrorMessage = err.Error()
	}

	sc.logger.StructuredPrint(logFields)
}
