package logger

type Level = string

const (
	PanicLevel Level = "panic"
	ErrorLevel Level = "error"
	WarnLevel  Level = "warn"
	InfoLevel  Level = "info"
	DebugLevel Level = "debug"
)
