package logging

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Logger interface {
	Info(appName string, caller string, message string)
	Error(appName string, caller string, message string)
	Warn(appName string, caller string, message string)
	Debug(appName string, caller string, message string)
	Fatal(appName string, caller string, message string)
}

type StructuredLog struct {
	AppName   string     `json:"app_name"`
	Caller    string     `json:"caller"`
	Message   string     `json:"message"`
	Timestamp *time.Time `json:"timestamp"`
}

func NewStructuredLog(appName string, caller, message string) *StructuredLog {
	now := time.Now()
	return &StructuredLog{
		AppName:   appName,
		Caller:    caller,
		Message:   message,
		Timestamp: &now,
	}
}

type LogrusHandler struct {
	Logrus *logrus.Logger
}

func (l *LogrusHandler) Info(
	appName string,
	caller string,
	message string,
) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":    caller,
			"message":   message,
			"app_name":  appName,
			"timestamp": time.Now(),
		},
	).Info()
}

func (l *LogrusHandler) Error(
	appName string,
	caller string,
	message string,
) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":    caller,
			"message":   message,
			"app_name":  appName,
			"timestamp": time.Now(),
		},
	).Error()
}

func (l *LogrusHandler) Warn(
	appName string,
	caller string,
	message string,
) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":    caller,
			"message":   message,
			"app_name":  appName,
			"timestamp": time.Now(),
		},
	).Warn()
}

func (l *LogrusHandler) Debug(
	appName string,
	caller string,
	message string,
) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":    caller,
			"message":   message,
			"app_name":  appName,
			"timestamp": time.Now(),
		},
	).Debug()
}

func (l *LogrusHandler) Fatal(
	appName string,
	caller string,
	message string,
) {
	l.Logrus.WithFields(
		logrus.Fields{
			"caller":    caller,
			"message":   message,
			"app_name":  appName,
			"timestamp": time.Now(),
		},
	).Fatal()
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
