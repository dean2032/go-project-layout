package logging

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global zap Logger with pid field
	zapLogger *zap.Logger
	// default sentry client
	sentryClient *sentry.Client
	// outPaths for zap log
	outPaths = []string{"stdout"}
	// initialFields use server ip, default value is pid
	initialFields = map[string]interface{}{
		"server_ip": ServerIP(),
	}
	// loggerName is default logger name
	loggerName = "logger"
	// atomicLevel is default logger atomic level
	atomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	// EncoderConfig default name config of log fields
	EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   CallerEncoder,
	}

	// AtomicLevelMap string level mapping zap AtomicLevel
	AtomicLevelMap = map[string]zap.AtomicLevel{
		"debug":  zap.NewAtomicLevelAt(zap.DebugLevel),
		"info":   zap.NewAtomicLevelAt(zap.InfoLevel),
		"warn":   zap.NewAtomicLevelAt(zap.WarnLevel),
		"error":  zap.NewAtomicLevelAt(zap.ErrorLevel),
		"dpanic": zap.NewAtomicLevelAt(zap.DPanicLevel),
		"panic":  zap.NewAtomicLevelAt(zap.PanicLevel),
		"fatal":  zap.NewAtomicLevelAt(zap.FatalLevel),
	}
	rwMutex sync.RWMutex
)

// Options new logger options
type Options struct {
	Name              string                 // logger name
	Level             string                 // logger level: debug, info, warn, error dpanic, panic, fatal
	Format            string                 // log format
	OutputPaths       []string               // log path
	InitialFields     map[string]interface{} // initial field names
	DisableCaller     bool                   // whether print caller
	DisableStacktrace bool                   // whether print stackstrace
	SentryClient      *sentry.Client         // sentry client
	EncoderConfig     *zapcore.EncoderConfig // logger key name config
	LumberjackSink    *LumberjackSink        // lumberjack sink for log rotate
}

// SentryDSNEnvKey ...
const (
	// SentryDSNEnvKey is sentry dsn environment variable, to init logger when import package
	SentryDSNEnvKey = "SENTRY_DSN"
	// SentryDebugEnvKey is environment variable to specified whether enable debug mode
	SentryDebugEnvKey = "SENTRY_DEBUG"
)

// init the global logger
func init() {
	var err error
	if dsn := os.Getenv(SentryDSNEnvKey); dsn != "" {
		debugStr := os.Getenv(SentryDebugEnvKey)
		debug := false
		if strings.ToLower(debugStr) != "" {
			debug = true
		}
		sentryClient, err = NewSentryClient(dsn, debug)
		if err != nil {
			log.Println(err)
		}
	}

	options := Options{
		Name:              loggerName,
		Level:             "debug",
		Format:            "json",
		OutputPaths:       outPaths,
		InitialFields:     initialFields,
		DisableCaller:     false,
		DisableStacktrace: true,
		SentryClient:      sentryClient,
		EncoderConfig:     &EncoderConfig,
		LumberjackSink:    nil,
	}
	zapLogger, err = NewLogger(options)
	if err != nil {
		log.Println(err)
	}
}

// InitLogger init logger
func InitLogger(name string, logDir string, debug bool) (*zap.Logger, error) {
	logPath := filepath.Join(logDir, "/access.log")
	if debug {
		outputPath := []string{"lumberjack:", "stderr"}
		sink := NewLumberjackSink("lumberjack", logPath, 30, 0, 0, false, true)
		return NewLogger(
			Options{
				Name:              name,
				Level:             "debug",
				Format:            "console",
				LumberjackSink:    sink,
				OutputPaths:       outputPath,
				DisableStacktrace: true,
			},
		)
	} else {
		outputPath := []string{"lumberjack:"}
		sink := NewLumberjackSink("lumberjack", logPath, 30, 10, 0, false, true)
		return NewLogger(Options{
			Name:              name,
			Level:             "info",
			Format:            "json",
			OutputPaths:       outputPath,
			LumberjackSink:    sink,
			DisableStacktrace: true,
		})
	}
}

// NewLogger return a zap Logger instance
func NewLogger(options Options) (*zap.Logger, error) {
	cfg := zap.Config{}
	// 设置日志级别
	lvl := strings.ToLower(options.Level)
	if _, exists := AtomicLevelMap[lvl]; !exists {
		cfg.Level = atomicLevel
	} else {
		cfg.Level = AtomicLevelMap[lvl]
		atomicLevel = cfg.Level
	}
	// use json by default for log format
	if strings.ToLower(options.Format) == "console" {
		cfg.Encoding = "console"
	} else {
		cfg.Encoding = "json"
	}
	// default output to stderr
	if len(options.OutputPaths) == 0 {
		cfg.OutputPaths = outPaths
		cfg.ErrorOutputPaths = outPaths
	} else {
		cfg.OutputPaths = options.OutputPaths
		cfg.ErrorOutputPaths = options.OutputPaths
	}
	if len(options.InitialFields) > 0 {
		for k, v := range options.InitialFields {
			initialFields[k] = v
		}
	}
	cfg.InitialFields = initialFields
	cfg.DisableCaller = options.DisableCaller
	cfg.DisableStacktrace = options.DisableStacktrace

	if options.EncoderConfig == nil {
		cfg.EncoderConfig = EncoderConfig
	} else {
		cfg.EncoderConfig = *options.EncoderConfig
	}

	if options.LumberjackSink != nil {
		if err := RegisterLumberjackSink(options.LumberjackSink); err != nil {
			Error("RegisterSink error", zap.Error(err))
		}
	}

	// build logger
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	if options.SentryClient != nil {
		logger = SentryAttach(logger, options.SentryClient)
	}

	if options.Name != "" {
		logger = logger.Named(options.Name)
	} else {
		logger = logger.Named(loggerName)
	}
	return logger, nil
}

// CloneLogger return the global logger copy which add a new name
func CloneLogger(name string, fields ...zap.Field) *zap.Logger {
	logger := zapLogger.Named(name)
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	return logger
}

// AttachCore add a core to zap logger
func AttachCore(l *zap.Logger, c zapcore.Core) *zap.Logger {
	return l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, c)
	}))
}

// ReplaceLogger replace global logger as newLogger
func ReplaceLogger(newLogger *zap.Logger) func() {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	prevLogger := zapLogger
	zapLogger = newLogger
	return func() { ReplaceLogger(prevLogger) }
}

// TextLevel get logger level
func TextLevel() string {
	b, _ := atomicLevel.MarshalText()
	return string(b)
}

// SetLevel set log level
func SetLevel(lvl string) {
	Warn("Set logger atomicLevel " + lvl)
	atomicLevel.UnmarshalText([]byte(strings.ToLower(lvl)))
}

// SentryClient get sentry client
func SentryClient() *sentry.Client {
	return sentryClient
}

// ServerIP get server IP
func ServerIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
