package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/segmentio/kafka-go"
)

func (h *Handler) publishMessageToKafka(ctx context.Context, parentMsg *model.PublishedKafkaMsg, traceID, messageType string, childMsg interface{}) error {
	ctx, span := h.tracer.Start(ctx, "Handler.publishMessageToKafka")
	defer span.End()

	msgBytes, err := json.Marshal(parentMsg)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Value: msgBytes,
		Time:  time.Now().UTC(),
	}

	carrier := tracerClient.NewKafkaHeadersCarrier(&kafkaMsg.Headers)
	tracerClient.InjectKafkaTracingHeadersToCarrier(ctx, carrier)

	injectTraceIDToKafkaHeaders(&kafkaMsg.Headers, traceID)

	if err := h.producerMap[messageType].PublishMessage(ctx, h.tracer, kafkaMsg); err != nil {
		h.logKafkaMessage(traceID, h.producerTopicMap[messageType], childMsg, err, "Error to publish message")
		return err
	}

	h.logKafkaMessage(traceID, h.producerTopicMap[messageType], childMsg, nil, "Success to publish message")
	return nil
}

func (h *Handler) logKafkaMessage(traceID, topic string, data interface{}, err error, activity string) {
	logFields := h.getKafkaLogFields(traceID, topic, data, err, activity)
	if err != nil {
		logFields.LogLevel = constants.LEVEL_ERROR
		logFields.Error = constants.TRUE
		logFields.Message.ErrorMessage = err.Error()
	}
	h.logger.StructuredPrint(logFields)
}

func (h *Handler) getKafkaLogFields(traceID, topic string, data interface{}, err error, activity string) *logger.LogFields {
	return &logger.LogFields{
		Timestamp:     time.Now().Format(constants.TIME_LAYOUT_FORMAT),
		LogLevel:      getLoggerLogLevel(err),
		TransactionID: "",
		ServiceName:   h.cfg.Project.ServiceName,
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
		ServerIP:          h.cfg.Project.ServerIP,
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

func (h *Handler) logRestMessage(traceID, endpoint, method string, data interface{}, err error, activity string) {
	logFields := h.getRestLogFields(traceID, endpoint, method, data, err, activity)
	if err != nil {
		logFields.LogLevel = constants.LEVEL_ERROR
		logFields.Error = constants.TRUE
		logFields.Message.ErrorMessage = err.Error()
	}
	h.logger.StructuredPrint(logFields)
}

func (h *Handler) getRestLogFields(traceID, endpoint, method string, data interface{}, err error, activity string) *logger.LogFields {
	return &logger.LogFields{
		Timestamp:     time.Now().Format(constants.TIME_LAYOUT_FORMAT),
		LogLevel:      getLoggerLogLevel(err),
		TransactionID: "",
		ServiceName:   h.cfg.Project.ServiceName,
		Endpoint:      endpoint,
		Protocol:      constants.PROTOCOL_HTTP,
		MethodType:    method,
		ExecutionType: constants.ASYNC,
		ContentType:   constants.CONTENT_JSON,
		FunctionName:  "",
		UserInfo: &logger.UserInfo{
			Username: "",
			Role:     "",
			Others:   "",
		},
		ExecutionTime:     "",
		ServerIP:          h.cfg.Project.ServerIP,
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

func injectTraceIDToKafkaHeaders(headers *[]kafka.Header, traceID string) {
	*headers = append(*headers, kafka.Header{
		Key:   "trace_id",
		Value: []byte(traceID),
	})
}
