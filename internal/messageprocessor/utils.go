package messageprocessor

import (
	"context"
	"fmt"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	"github.com/segmentio/kafka-go"
)

func (mp *MessageProcessor) getProcessorFunc(categoryName string) (func(context.Context, *kafka.Reader, kafka.Message, *model.ConsumedKafkaMsg, string, int, string), bool) {
	processorFuncMap := map[string]func(context.Context, *kafka.Reader, kafka.Message, *model.ConsumedKafkaMsg, string, int, string){
		constants.NOTIF_TYPE_EMAIL: mp.processEmail,
		constants.NOTIF_TYPE_SMS:   mp.processSms,
		constants.NOTIF_TYPE_INAPP: mp.processInApp,
		constants.NOTIF_TYPE_PUSH:  mp.processPush,
	}

	processorFunc, exists := processorFuncMap[categoryName]
	return processorFunc, exists
}

func (mp *MessageProcessor) commitAndLogMsg(ctx context.Context, r *kafka.Reader, kafkaMsg kafka.Message, childMsg interface{}, traceID string) {
	if err := r.CommitMessages(ctx, kafkaMsg); err != nil {
		mp.logKafkaMessage(traceID, kafkaMsg.Topic, childMsg, err, "Error while committing message")
	}
	mp.logKafkaMessage(traceID, kafkaMsg.Topic, childMsg, nil, "Success to commit kafka message")
}

func (mp *MessageProcessor) logKafkaMessage(traceID, topic string, data interface{}, err error, activity string) {
	logFields := mp.getKafkaLogFields(traceID, topic, data, err, activity)
	if err != nil {
		logFields.LogLevel = constants.LEVEL_ERROR
		logFields.Error = constants.TRUE
		logFields.Message.ErrorMessage = err.Error()
	}
	mp.logger.StructuredPrint(logFields)
}

func (mp *MessageProcessor) getKafkaLogFields(traceID, topic string, data interface{}, err error, activity string) *logger.LogFields {
	return &logger.LogFields{
		Timestamp:     time.Now().Format(constants.TIME_LAYOUT_FORMAT),
		LogLevel:      getLoggerLogLevel(err),
		TransactionID: "",
		ServiceName:   mp.cfg.Project.ServiceName,
		Endpoint:      topic,
		Protocol:      constants.PROTOCOL_TCP,
		MethodType:    constants.KAFKA_WRITER,
		ExecutionType: constants.ASYNC,
		ContentType:   constants.CONTENT_JSON,
		FunctionName:  "",
		UserInfo: &logger.UserInfo{
			Username: "",
			Role:     "",
			Others:   "",
		},
		ExecutionTime:     "",
		ServerIP:          mp.cfg.Project.ServerIP,
		ClientIP:          "",
		EventName:         "",
		TraceID:           traceID,
		PrevTransactionID: "",
		Body:              fmt.Sprint(data),
		Result:            "",
		Error:             constants.FALSE,
		FlagStartOrStop:   constants.STOP,
		Message: &logger.Message{
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

func createKafkaHeadersMap(headers []kafka.Header) map[string]string {
	headersMap := make(map[string]string)
	for _, header := range headers {
		headersMap[header.Key] = string(header.Value)
	}
	return headersMap
}

func getValueFromKafkaHeaders(headers map[string]string, key string) string {
	return headers[key]
}

func createPulishedKafkaMessage(consumedKafkaMsg *model.ConsumedKafkaMsg) *model.PublishedKafkaMsg {
	return &model.PublishedKafkaMsg{
		TypeName:      consumedKafkaMsg.TypeName,
		CategoryName:  consumedKafkaMsg.CategoryName,
		ChannelName:   consumedKafkaMsg.ChannelName,
		PriorityOrder: consumedKafkaMsg.PriorityOrder,
		PreferenceUrl: consumedKafkaMsg.PreferenceUrl,
		Data:          consumedKafkaMsg.Data,
	}
}
