package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// CtxLoggerName ...
const (
	// CtxLoggerName define the ctx logger name
	CtxLoggerName = "ctx"
	// TraceIDKeyName define the trace id keyname
	TraceIDKeyName = "trace_id"
	// TraceIDHeaderName ...
	TraceIDHeaderName = "X-Trace-Id"
)

func getLoggerFromCtx(c context.Context) *zap.Logger {
	if c == nil {
		return nil
	}
	var ctxLoggerItf interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLoggerItf, _ = gc.Get(CtxLoggerName)
	} else {
		ctxLoggerItf = c.Value(CtxLoggerName)
	}
	if ctxLoggerItf != nil {
		return ctxLoggerItf.(*zap.Logger)
	}
	return nil
}

// CtxLogger get the ctxLogger in context
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	if c == nil {
		c = context.Background()
	}
	ctxLogger := getLoggerFromCtx(c)
	if ctxLogger == nil {
		_, ctxLogger = NewCtxLogger(c, CloneLogger(CtxLoggerName), CtxTraceID(c))
	}
	if len(fields) > 0 {
		ctxLogger = ctxLogger.With(fields...)
	}
	return ctxLogger
}

// CtxTraceID get trace id from context
// Modify TraceIDPrefix change change the prefix
func CtxTraceID(c context.Context) string {
	if c == nil {
		c = context.Background()
	}
	if gc, ok := c.(*gin.Context); ok {
		if traceID := gc.GetString(TraceIDKeyName); traceID != "" {
			return traceID
		}
	}
	traceIDItf := c.Value(TraceIDKeyName)
	if traceIDItf != nil {
		return traceIDItf.(string)
	}
	// return default value
	return xid.New().String()
}

// NewCtxLogger return a context with logger and trace id and a logger with trace id
func NewCtxLogger(c context.Context, logger *zap.Logger, traceID string) (context.Context, *zap.Logger) {
	if c == nil {
		c = context.Background()
	}
	if traceID == "" {
		traceID = CtxTraceID(c)
	}
	ctxLogger := logger.With(zap.String(TraceIDKeyName, traceID))
	if gc, ok := c.(*gin.Context); ok {
		gc.Set(CtxLoggerName, ctxLogger)
		gc.Set(TraceIDKeyName, traceID)
	} else {
		c = context.WithValue(c, TraceIDKeyName, traceID)
		c = context.WithValue(c, CtxLoggerName, ctxLogger)
	}
	return c, ctxLogger
}
