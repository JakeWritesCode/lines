package internal

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Logger interface {
	Info(log *StructuredLog)
	Error(log *StructuredLog)
	Warn(log *StructuredLog)
	Debug(log *StructuredLog)
	Fatal(log *StructuredLog)
	EmitLog(level string, caller string, message string)
}

type StructuredLog struct {
	Caller    string     `json:"caller"`
	Message   string     `json:"message"`
	Timestamp *time.Time `json:"timestamp"`
}

func NewStructuredLog(caller, message string) *StructuredLog {
	now := time.Now()
	return &StructuredLog{
		Caller:    caller,
		Message:   message,
		Timestamp: &now,
	}
}

type LogrusHandler struct {
	Logrus *logrus.Logger
}

func (l *LogrusHandler) Info(log *StructuredLog) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":  log.Caller,
			"message": log.Message,
		},
	).Info()
}

func (l *LogrusHandler) Error(log *StructuredLog) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":  log.Caller,
			"message": log.Message,
		},
	).Error()
}

func (l *LogrusHandler) Warn(log *StructuredLog) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":  log.Caller,
			"message": log.Message,
		},
	).Warn()
}

func (l *LogrusHandler) Debug(log *StructuredLog) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":  log.Caller,
			"message": log.Message,
		},
	).Debug()
}

func (l *LogrusHandler) Fatal(log *StructuredLog) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":  log.Caller,
			"message": log.Message,
		},
	).Fatal()
}

func (l *LogrusHandler) EmitLog(level string, caller string, message string) {
	structuredLog := NewStructuredLog(caller, message)
	switch level {
	case "info":
		l.Info(structuredLog)
	case "error":
		l.Error(structuredLog)
	case "warn":
		l.Warn(structuredLog)
	case "debug":
		l.Debug(structuredLog)
	case "fatal":
		l.Fatal(structuredLog)
	default:
		l.Info(structuredLog)
	}
}

func NewLogrusHandler(level string) *LogrusHandler {
	logrusInstance := logrus.New()
	logrusInstance.SetFormatter(&logrus.JSONFormatter{})
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrusInstance.SetLevel(logrus.InfoLevel)
	} else {
		logrusInstance.SetLevel(parsedLevel)
	}
	logger := LogrusHandler{
		Logrus: logrusInstance,
	}
	return &logger
}
