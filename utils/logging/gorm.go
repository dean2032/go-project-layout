// gorm v2

package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLoggerName ...
var (
	// GormLoggerName is name of gorm logger
	GormLoggerName = "gorm"
	// GormLoggerCallerSkip caller skip
	GormLoggerCallerSkip = 3
)

// GormLogger use zap to print gorm log
type GormLogger struct {
	logLevel      zapcore.Level
	slowThreshold time.Duration
}

var gormLogLevelMap = map[gormlogger.LogLevel]zapcore.Level{
	gormlogger.Info:  zap.InfoLevel,
	gormlogger.Warn:  zap.WarnLevel,
	gormlogger.Error: zap.ErrorLevel,
}

// LogMode implements gorm logger interface method
func (g GormLogger) LogMode(gormLogLevel gormlogger.LogLevel) gormlogger.Interface {
	zaplevel, exists := gormLogLevelMap[gormLogLevel]
	if !exists {
		zaplevel = zap.DebugLevel
	}
	newlogger := g
	newlogger.logLevel = zaplevel
	return &newlogger
}

// CtxLogger create ctxlogger
func (g GormLogger) CtxLogger(ctx context.Context) *zap.Logger {
	return CtxLogger(ctx).Named(GormLoggerName).WithOptions(zap.AddCallerSkip(GormLoggerCallerSkip))
}

// Info implements gorm logger interface method
func (g GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.InfoLevel {
		g.CtxLogger(ctx).Sugar().Infof(msg, data...)
	}
}

// Warn implements gorm logger interface method
func (g GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.WarnLevel {
		g.CtxLogger(ctx).Sugar().Warnf(msg, data...)
	}
}

// Error implements gorm logger interface method
func (g GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.ErrorLevel {
		g.CtxLogger(ctx).Sugar().Errorf(msg, data...)
	}
}

// Trace implements gorm logger interface method
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	now := time.Now()
	latency := now.Sub(begin).Seconds()
	switch {
	case err != nil:
		sql, rows := fc()
		if err != gorm.ErrRecordNotFound {
			g.CtxLogger(
				ctx,
			).Error(
				"sql: "+sql,
				zap.Float64("latency", latency),
				zap.Int64("rows", rows),
				zap.String("error", err.Error()),
			)
		}
	case g.slowThreshold != 0 && latency > g.slowThreshold.Seconds():
		sql, rows := fc()
		g.CtxLogger(
			ctx,
		).Warn(
			"sql: "+sql,
			zap.Float64("latency", latency),
			zap.Int64("rows", rows),
			zap.Float64("threshold", g.slowThreshold.Seconds()),
		)
	case g.logLevel <= zap.InfoLevel:
		sql, rows := fc()
		g.CtxLogger(ctx).Info("sql: "+sql, zap.Float64("latency", latency), zap.Int64("rows", rows))
	}
}

// NewGormLogger make GormLogger
func NewGormLogger(logLevel zapcore.Level, slowThreshold time.Duration) GormLogger {
	return GormLogger{
		logLevel:      logLevel,
		slowThreshold: slowThreshold,
	}
}
