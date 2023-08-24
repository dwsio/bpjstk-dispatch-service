package constants

const (
	TIME_LAYOUT_FORMAT = "2006-01-02 15:04:05"

	LEVEL_INFO  = "INFO"
	LEVEL_WARN  = "WARN"
	LEVEL_ERROR = "ERROR"
	LEVEL_FATAL = "FATAL"

	METHOD_POST = "POST"
	METHOD_GET  = "GET"

	KAFKA_READER = "KAFKA_READER"
	KAFKA_WRITER = "KAFKA_WRITER"

	PROTOCOL_HTTP = "HTTP"
	PROTOCOL_TCP  = "TCP"

	SYNC  = "SYNC"
	ASYNC = "ASYNC"

	TRUE  = "TRUE"
	FALSE = "FALSE"

	START = "START"
	STOP  = "STOP"

	CONTENT_JSON = "JSON"

	PRIORITY_HIGH   = "high"
	PRIORITY_NORMAL = "reg"

	StatusAborted          = "Message processing aborted"
	ErrorProcessingMessage = "Error while processing message"

	Success = "Sukses"
	Failed  = "Gagal"

	NOTIF_TYPE_EMAIL    = "email"
	NOTIF_TYPE_SMS      = "sms"
	NOTIF_TYPE_INAPP    = "inapp"
	NOTIF_TYPE_PUSH     = "push"
	NOTIF_TYPE_SMS_POOL = "sms_pool"

	CHANNEL_JMO     = "jmo"
	CHANNEL_SMILE   = "smile"
	CHANNEL_SIPP    = "sipp"
	CHANNEL_SIDIA   = "sidia"
	CHANNEL_PERISAI = "perisai"
)
