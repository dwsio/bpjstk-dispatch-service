package logger

import (
	"log"
	"os"
	"reflect"
)

type Config struct {
	Encoding string
}

type LogFields struct {
	Timestamp         string
	LogLevel          string
	TransactionID     string
	ServiceName       string
	Endpoint          string
	Protocol          string
	MethodType        string
	ExecutionType     string
	ContentType       string
	FunctionName      string
	UserInfo          *UserInfo
	ExecutionTime     string
	ServerIP          string
	ClientIP          string
	EventName         string
	TraceID           string
	PrevTransactionID string
	Body              string
	Result            string
	Error             string
	FlagStartOrStop   string
	Message           *Message
}

type UserInfo struct {
	Username string
	Role     string
	Others   string
}

type Message struct {
	Activity          string
	ObjectPerformedOn string
	ResultOfActivity  string
	ErrorCode         string
	ErrorMessage      string
	ShortDescription  string
}

type AppLogger struct {
	logger *log.Logger
}

func NewAppLogger() *AppLogger {
	return &AppLogger{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (al *AppLogger) StructuredPrint(lf *LogFields) {
	wrapEmptyFields(lf)

	logString := lf.Timestamp + " [" + lf.LogLevel + "] \t " + lf.TransactionID + " \t " + lf.ServiceName + " \t " + lf.Endpoint + " \t " + lf.Protocol + " \t " + lf.MethodType + " \t " + lf.ExecutionType + " \t " + lf.ContentType + " \t " + lf.FunctionName + " \t '" + lf.UserInfo.Username + "' as '" + lf.UserInfo.Role + "' . '" + lf.UserInfo.Others + "' \t " + lf.ExecutionTime + " ms \t " + lf.ServerIP + " \t " + lf.ClientIP + " \t " + lf.EventName + " \t " + lf.TraceID + " \t " + lf.PrevTransactionID + " \t " + lf.Body + " \t " + lf.Result + " \t " + lf.Error + " \t [" + lf.FlagStartOrStop + "] \t '" + lf.Message.Activity + "' on '" + lf.Message.ObjectPerformedOn + "' with result '" + lf.Message.ShortDescription + "' with error '" + lf.Message.ErrorMessage + "' : '" + lf.Message.ErrorCode + "' . '" + lf.Message.ShortDescription + "'"

	al.logger.Println(logString)
}

func wrapEmptyFields(fields interface{}) {
	v := reflect.ValueOf(fields).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			wrapEmptyFields(field.Interface())
		} else if field.Kind() == reflect.Struct {
			wrapEmptyFields(field.Addr().Interface())
		} else if field.Kind() == reflect.String && field.String() == "" {
			field.SetString("-")
		}
	}
}
