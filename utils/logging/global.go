package logging

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// Debug ...
func Debug(msg string, fields ...zap.Field) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Info ...
func Info(msg string, fields ...zap.Field) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...zap.Field) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...zap.Field) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...zap.Field) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).DPanic(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...zap.Field) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...zap.Field) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// Debugf ...
func Debugf(template string, args ...interface{}) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar().Debugf(template, args...)
}

// Infof ...
func Infof(template string, args ...interface{}) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar().Infof(template, args...)
}

// Warnf ...
func Warnf(template string, args ...interface{}) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar().Warnf(template, args...)
}

// Errorf ...
func Errorf(template string, args ...interface{}) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar().Errorf(template, args...)
}

// DPanicf ...
func DPanicf(template string, args ...interface{}) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar().DPanicf(template, args...)
}

// Panicf ...
func Panicf(template string, args ...interface{}) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar().Panicf(template, args...)
}

// Fatalf ...
func Fatalf(template string, args ...interface{}) {
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar().Fatalf(template, args...)
}

// SentryCaptureMessage ...
func SentryCaptureMessage(msg string) error {
	if SentryClient() == nil {
		return errors.New("sentry client is nil, please set the sentry dsn config")
	}
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage(msg)
	return nil
}

// SentryCaptureException ...
func SentryCaptureException(err error) error {
	if SentryClient() == nil {
		return errors.New("sentry client is nil, please set the sentry dsn config")
	}
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureException(err)
	return nil
}
