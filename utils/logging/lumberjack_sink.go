package logging

import (
	"net/url"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

// LogFilename ...
const (
	// LogFilename default logger file name
	LogFilename = "/tmp/logger.log"
)

// LumberjackSink for log rotate
type LumberjackSink struct {
	*lumberjack.Logger
	Scheme string
}

// Sync lumberjack Logger implements Sync method for Sink interface
func (LumberjackSink) Sync() error {
	return nil
}

// RegisterLumberjackSink register lumberjack sink
func RegisterLumberjackSink(sink *LumberjackSink) error {
	err := zap.RegisterSink(sink.Scheme, func(*url.URL) (zap.Sink, error) {
		if sink.Filename == "" {
			sink.Filename = LogFilename
		}
		return sink, nil
	})
	return err
}

// NewLumberjackSink build LumberjackSink
func NewLumberjackSink(
	scheme, filename string,
	maxAge, maxBackups, maxSize int,
	compress, localtime bool,
) *LumberjackSink {
	return &LumberjackSink{
		Logger: &lumberjack.Logger{
			Filename:   filename,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
			MaxSize:    maxSize,
			Compress:   compress,
			LocalTime:  localtime,
		},
		Scheme: scheme,
	}
}
