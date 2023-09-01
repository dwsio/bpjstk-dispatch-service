package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	contextMd "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/context_metadata"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/logger"
	tracerClient "git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/tracer"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/metric"
)

func (u *Usecase) publishMessageToKafka(ctx context.Context, parentMsg *model.PublishedKafkaMsg, messageType string, childMsg interface{}) error {
	ctx, span := u.tracer.Start(ctx, "Usecase.publishMessageToKafka")
	defer span.End()

	md, _ := contextMd.GetMetadataFromContext(ctx)

	msgBytes, err := json.Marshal(parentMsg)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Value: msgBytes,
		Time:  time.Now().UTC(),
	}

	tracerClient.InjectKafkaTracingHeadersToCarrier(ctx, &kafkaMsg.Headers)

	injectTraceIDToKafkaHeaders(&kafkaMsg.Headers, md.TraceID)

	if err := u.producerMap[messageType].PublishMessage(ctx, kafkaMsg); err != nil {
		u.serviceMetrics.ErrorKafkaPublish.Add(context.Background(), 1, metric.WithAttributes())
		u.logKafkaMessage(ctx, childMsg, err, "Error to publish message")
		return err
	}

	u.serviceMetrics.SuccessKafkaPublish.Add(context.Background(), 1, metric.WithAttributes())
	u.logKafkaMessage(ctx, childMsg, nil, "Success to publish message")
	return nil
}

func (u *Usecase) logKafkaMessage(ctx context.Context, data interface{}, err error, activity string) {

	logFields := u.getKafkaLogFields(ctx, data, err, activity)
	if err != nil {
		logFields.LogLevel = constants.LEVEL_ERROR
		logFields.Error = constants.TRUE
		logFields.Message.ErrorMessage = err.Error()
	}
	u.logger.StructuredPrint(logFields)
}

func (u *Usecase) getKafkaLogFields(ctx context.Context, data interface{}, err error, activity string) *logger.LogFields {
	metadata, _ := contextMd.GetMetadataFromContext(ctx)

	return &logger.LogFields{
		Timestamp:     time.Now().Format(constants.TIME_LAYOUT_FORMAT),
		LogLevel:      getLoggerLogLevel(err),
		TransactionID: fmt.Sprintf("TR%s", metadata.TransactionID),
		ServiceName:   u.cfg.Project.ServiceName,
		Endpoint:      metadata.Topic,
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
		ExecutionTime:     fmt.Sprint(time.Since(time.UnixMilli(metadata.StartTime)).Milliseconds()),
		ServerIP:          u.cfg.Project.ServerIP,
		ClientIP:          "",
		EventName:         "",
		TraceID:           fmt.Sprintf("TC%s", metadata.TraceID),
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

func (u *Usecase) logRestMessage(ctx context.Context, endpoint, method string, data interface{}, err error, activity string) {
	logFields := u.getRestLogFields(ctx, endpoint, method, data, err, activity)
	if err != nil {
		logFields.LogLevel = constants.LEVEL_ERROR
		logFields.Error = constants.TRUE
		logFields.Message.ErrorMessage = err.Error()
	}
	u.logger.StructuredPrint(logFields)
}

func (u *Usecase) getRestLogFields(ctx context.Context, endpoint, method string, data interface{}, err error, activity string) *logger.LogFields {
	metadata, _ := contextMd.GetMetadataFromContext(ctx)

	return &logger.LogFields{
		Timestamp:     time.Now().Format(constants.TIME_LAYOUT_FORMAT),
		LogLevel:      getLoggerLogLevel(err),
		TransactionID: fmt.Sprintf("TR%s", metadata.TransactionID),
		ServiceName:   u.cfg.Project.ServiceName,
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
		ExecutionTime:     fmt.Sprint(time.Since(time.UnixMilli(metadata.StartTime)).Milliseconds()),
		ServerIP:          u.cfg.Project.ServerIP,
		ClientIP:          "",
		EventName:         "",
		TraceID:           fmt.Sprintf("TC%s", metadata.TraceID),
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
