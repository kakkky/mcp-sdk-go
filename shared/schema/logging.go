package schema

type LoggingLevelSchema string

const (
	DEBUG     LoggingLevelSchema = "debug"
	INFO      LoggingLevelSchema = "info"
	NOTICE    LoggingLevelSchema = "notice"
	WARNING   LoggingLevelSchema = "warning"
	ERROR     LoggingLevelSchema = "error"
	CRITICAL  LoggingLevelSchema = "critical"
	ALERT     LoggingLevelSchema = "alert"
	EMERGENCY LoggingLevelSchema = "emergency"
)
