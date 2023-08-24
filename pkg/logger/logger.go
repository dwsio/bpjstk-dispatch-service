package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// LogLevelDebug for debugging information.
	LogLevelDebug string = "debug"
	// LogLevelInfo for informational messages.
	LogLevelInfo string = "info"
	// LogLevelError for error messages.
	LogLevelError string = "error"
)

// LoggerConfig holds the configuration for creating a new logger.
type Config struct {
	Encoding   string
	Level      string
	OutputPath string
	ErrorPath  string
}

// AppLogger provides structured logging functionality.
type AppLogger struct {
	Logger  *zap.Logger
	SLogger *zap.SugaredLogger
}

type LogFields struct {
	Timestamp         string    `json:"timestamp"`
	LogLevel          string    `json:"log_level"`
	TransactionID     string    `json:"transaction_id"`
	ServiceName       string    `json:"service_name"`
	Endpoint          string    `json:"endpoint"`
	Protocol          string    `json:"protocol"`
	MethodType        string    `json:"method_type"`
	ExecutionType     string    `json:"execution_type"`
	ContentType       string    `json:"content_type"`
	FunctionName      string    `json:"function_name"`
	UserInfo          *UserInfo `json:"user_info"`
	ExecutionTime     string    `json:"execution_time"`
	ServerIP          string    `json:"server_ip"`
	ClientIP          string    `json:"client_ip"`
	EventName         string    `json:"event_name"`
	TraceID           string    `json:"trace_id"`
	PrevTransactionID string    `json:"prev_transaction_id"`
	Body              string    `json:"body"`
	Result            string    `json:"result"`
	Error             string    `json:"error"`
	FlagStartOrStop   string    `json:"flag_start_or_stop"`
	Message           *Message  `json:"message"`
}

type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Others   string `json:"others"`
}

type Message struct {
	Activity          string `json:"activity"`
	ObjectPerformedOn string `json:"object_performed_on"`
	ResultOfActivity  string `json:"result_of_activity"`
	ErrorCode         string `json:"error_code"`
	ErrorMessage      string `json:"error_message"`
	ShortDescription  string `json:"short_description"`
}

// NewAppLogger creates a new AppLogger instance.
func NewAppLogger(cfg *Config) (*AppLogger, error) {
	var zapLevel zapcore.Level
	switch cfg.Level {
	case LogLevelDebug:
		zapLevel = zapcore.DebugLevel
	case LogLevelError:
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	cfgZap := zap.Config{
		Encoding:         cfg.Encoding,
		Level:            zap.NewAtomicLevelAt(zapLevel),
		OutputPaths:      []string{cfg.OutputPath},
		ErrorOutputPaths: []string{cfg.ErrorPath},
	}

	logger, err := cfgZap.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return &AppLogger{
		Logger:  logger,
		SLogger: logger.Sugar(),
	}, nil
}

// StructuredPrint logs the given LogFields.
func (al *AppLogger) StructuredPrint(lf *LogFields) {
	if lf == nil {
		al.SLogger.Error("Print called with nil LogFields")
		return
	}

	al.Logger.Info("",
		zap.String("timestamp", lf.Timestamp),
		zap.String("log_level", lf.LogLevel),
		zap.String("transaction_id", lf.TransactionID),
		zap.String("service_name", lf.ServiceName),
		zap.String("endpoint", lf.Endpoint),
		zap.String("protocol", lf.Protocol),
		zap.String("method_type", lf.MethodType),
		zap.String("execution_type", lf.ExecutionType),
		zap.String("content_type", lf.ContentType),
		zap.String("function_name", lf.FunctionName),
		zap.Any("user_info", lf.UserInfo),
		zap.String("execution_time", lf.ExecutionTime),
		zap.String("server_ip", lf.ServerIP),
		zap.String("client_ip", lf.ClientIP),
		zap.String("event_name", lf.EventName),
		zap.String("trace_id", lf.TraceID),
		zap.String("prev_transaction_id", lf.PrevTransactionID),
		zap.String("body", lf.Body),
		zap.String("result", lf.Result),
		zap.String("error", lf.Error),
		zap.String("flag_start_or_stop", lf.FlagStartOrStop),
		zap.Any("message", lf.Message),
	)
}

// LogError logs the error along with additional context.
func (al *AppLogger) LogError(err error, message string) {
	if err == nil {
		al.SLogger.Error("LogError called with nil error")
		return
	}

	al.Logger.Error(message,
		zap.Error(err),
	)
}
