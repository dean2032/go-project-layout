package logging

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	if zapLogger == nil {
		t.Error("logger is nil")
	}
}

func TestNewLoggerNoParam(t *testing.T) {
	logger, err := NewLogger(Options{})
	if err != nil {
		t.Error(err)
	}
	if logger == nil {
		t.Error("return a nil logger")
	}
	logger.Debug("TestNewLoggerNoParam Debug")
}

func TestCloneLogger(t *testing.T) {
	nlogger := CloneLogger("cloned")
	if reflect.DeepEqual(nlogger, zapLogger) {
		t.Error("CloneLogger should not be default logger")
	}
	if &nlogger == &zapLogger {
		t.Error("CloneLogger should not be default logger")
	}
}

func TestSetLevel(t *testing.T) {
	zapLogger.Debug("TestChangeLevel raw debug level")
	t.Log("current level:", atomicLevel.Level())
	atomicLevel.SetLevel(zap.InfoLevel)
	t.Log("new level:", atomicLevel.Level())
	zapLogger.Debug("TestChangeLevel raw debug level should not be logged")
	// reset
	atomicLevel.SetLevel(zap.DebugLevel)
}

func TestTextLevel(t *testing.T) {
	level := TextLevel()
	if level != "debug" {
		t.Error(level)
	}

}
