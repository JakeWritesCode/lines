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
	testLog := NewStructuredLog("test", "test message")
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

	testLog := NewStructuredLog("TestInfo", "This is a test log")
	logrusHandler.Info(testLog)

	logOutput := buf.String()
	expected := `{"caller":"TestInfo","level":"info","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Error(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	testLog := NewStructuredLog("TestError", "This is a test log")
	logrusHandler.Error(testLog)

	logOutput := buf.String()
	expected := `{"caller":"TestError","level":"error","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Warn(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	testLog := NewStructuredLog("TestWarn", "This is a test log")
	logrusHandler.Warn(testLog)

	logOutput := buf.String()
	expected := `{"caller":"TestWarn","level":"warning","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Debug(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	testLog := NewStructuredLog("TestDebug", "This is a test log")
	logrusHandler.Debug(testLog)

	logOutput := buf.String()
	expected := `{"caller":"TestDebug","level":"debug","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_Fatal(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf
	logrusHandler.Logrus.ExitFunc = func(int) {}

	testLog := NewStructuredLog("TestFatal", "This is a test log")
	logrusHandler.Fatal(testLog)

	logOutput := buf.String()
	expected := `{"caller":"TestFatal","level":"fatal","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_EmitLog_Info(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.EmitLog("info", "TestInfo", "This is a test log")

	logOutput := buf.String()
	expected := `{"caller":"TestInfo","level":"info","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_EmitLog_Error(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.EmitLog("error", "TestError", "This is a test log")

	logOutput := buf.String()
	expected := `{"caller":"TestError","level":"error","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_EmitLog_Warn(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.EmitLog("warn", "TestWarn", "This is a test log")

	logOutput := buf.String()
	expected := `{"caller":"TestWarn","level":"warning","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_EmitLog_Debug(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.EmitLog("debug", "TestDebug", "This is a test log")

	logOutput := buf.String()
	expected := `{"caller":"TestDebug","level":"debug","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}

func TestLogrusHandler_EmitLog_Default(t *testing.T) {
	var buf bytes.Buffer
	logrusHandler := NewLogrusHandler("debug")
	logrusHandler.Logrus.Out = &buf

	logrusHandler.EmitLog("invalid", "TestDefault", "This is a test log")

	logOutput := buf.String()
	expected := `{"caller":"TestDefault","level":"info","message":"This is a test log","msg":"","time":"`
	if len(logOutput) < len(expected) || logOutput[:len(expected)] != expected {
		t.Errorf("\n  output:  %s\nexpected: %s\n", logOutput, expected)
	}
}
