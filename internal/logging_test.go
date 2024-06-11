package internal

import (
	"bytes"
	"testing"
)

func TestNewLogrusHandler(t *testing.T) {
	logger := NewLogrusHandler("info")
	if logger == nil {
		t.Error("Expected logger to be created, but it was nil")
	}
	if logger.Logrus.Level != 4 {
		t.Error("Expected logger level to be info, but it was not")
	}
	if logger.Logrus.Formatter == nil {
		t.Error("Expected logger formatter to be set, but it was not")
	}
}

func TestNewLogrusHandler_FailedLevelParse(t *testing.T) {
	logger := NewLogrusHandler("invalid")
	if logger == nil {
		t.Error("Expected logger to be created, but it was nil")
	}
	if logger.Logrus.Level != 4 {
		t.Error("Expected logger level to be info, but it was not")
	}
	if logger.Logrus.Formatter == nil {
		t.Error("Expected logger formatter to be set, but it was not")
	}
}

func TestNewStructuredLog(t *testing.T) {
	testLog := NewStructuredLog("app", "test", "test message")
	if testLog.AppName != "app" {
		t.Error("Expected test log app name to be set, but it was not")
	}
	if testLog.Caller != "test" {
		t.Error("Expected test log caller to be set, but it was not")
	}
	if testLog.Message != "test message" {
		t.Error("Expected test log message to be set, but it was not")
	}
	if testLog.Timestamp == nil {
		t.Error("Expected test log timestamp to be set, but it was not")
	}
}

func TestLogrusHandler_Info(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.Info(
		"test",
		"TestInfo",
		"This is a test log",
	)

	logOutput := buf.String()
	expected := `{"app_name":"test","caller":"TestInfo","level":"info","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Error(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.Error("TestApp", "TestError", "This is a test log")

	logOutput := buf.String()
	expected := `{"app_name":"TestApp","caller":"TestError","level":"error","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Warn(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.Warn("TestApp", "TestWarn", "This is a test log")

	logOutput := buf.String()
	expected := `{"app_name":"TestApp","caller":"TestWarn","level":"warning","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Debug(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.Debug("TestApp", "TestDebug", "This is a test log")

	logOutput := buf.String()
	expected := `{"app_name":"TestApp","caller":"TestDebug","level":"debug","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Fatal(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf
	logrusHandler.Logrus.ExitFunc = func(int) {}

	logrusHandler.Fatal("TestApp", "TestFatal", "This is a test log")

	logOutput := buf.String()
	expected := `{"app_name":"TestApp","caller":"TestFatal","level":"fatal","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}
