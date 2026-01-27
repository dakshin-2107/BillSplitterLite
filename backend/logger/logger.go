package logger

import (
	"os"

	"go.uber.org/zap"
)

const env_true = "true"

type ILogger interface {
	Init()
	DebugLog(message string)
	InfoLog(message string)
	PanicLog(message string)
}

func NewLogger() ILogger {
	return &Logger{}
}

type Logger struct {
	allowDebugLogs bool
	allowiInfoLogs bool
	allowPanicLogs bool
	logger         *zap.Logger
}

func (l *Logger) Init() {

	temp, err := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1)) // AddCallerSkip(1) skips the logger function itself
	if err != nil {
		panic(err)
	}

	l.allowDebugLogs = os.Getenv("ALLOW_DEBUG_LOGS") == env_true
	l.allowiInfoLogs = os.Getenv("ALLOW_INFO_LOGS") == env_true
	l.allowPanicLogs = os.Getenv("ALLOW_PANIC_LOGS") == env_true
	l.logger = temp
}

// log about code flow
func (l *Logger) DebugLog(msg string) {
	if l.allowDebugLogs {
		l.logger.Log(zap.DebugLevel, msg)
	}
}

// log about flow of business logic
func (l *Logger) InfoLog(msg string) {
	if l.allowiInfoLogs {
		l.logger.Log(zap.InfoLevel, msg)
	}
}

// log something before call to panic
func (l *Logger) PanicLog(msg string) {
	if !l.allowPanicLogs {
		l.logger.Log(zap.WarnLevel, msg)
	}
}
