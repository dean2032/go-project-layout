package logging

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultGinSlowThreshold = time.Second * 3
)

// GetGinTraceIDFromHeader get key in request header of gin use for TraceIDKeyName as traceid
func GetGinTraceIDFromHeader(c *gin.Context) string {
	return c.Request.Header.Get(TraceIDHeaderName)
}

// GetGinTraceIDFromQueryString get key in query string of gin use for TraceIDKeyName as traceid
func GetGinTraceIDFromQueryString(c *gin.Context) string {
	return c.Query(TraceIDKeyName)
}

// GinLogDetails gin log details used by middleware
type GinLogDetails struct {
	Timestamp     time.Time `json:"timestamp"`
	Method        string    `json:"method"`
	Path          string    `json:"path"`
	Query         string    `json:"query"`
	Proto         string    `json:"proto"`
	ContentLength int       `json:"content_length"`
	Host          string    `json:"host"`
	RemoteAddr    string    `json:"remote_addr"`
	RequestURI    string    `json:"request_uri"`
	Referer       string    `json:"referer"`
	UserAgent     string    `json:"user_agent"`
	ClientIP      string    `json:"client_ip"`
	ContentType   string    `json:"content_type"`
	HandlerName   string    `json:"handler_name"`
	StatusCode    int       `json:"status_code"`
	BodySize      int       `json:"body_size"`
	Latency       float64   `json:"latency"`
}

// GinLoggerConfig define GinLogger fields
type GinLoggerConfig struct {
	// Optional. Default value is logger.defaultGinLogFormatter
	Formatter func(GinLogDetails) string
	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
	// SkipPathRegexps skip path by regexp
	SkipPathRegexps []string
	// TraceIDFunc function to get or generate trace id
	// Optional.
	TraceIDFunc func(context.Context) string
	// InitFieldsFunc function to init logger fields,  key is field name, value is field value
	InitFieldsFunc func(context.Context) map[string]interface{}
	// Whether print details
	// Optional.
	EnableDetails bool
	// slow request threshold, output Error log if request handle time exeeded this value
	SlowThreshold time.Duration

	skipPathMap map[string]struct{}
	skipRegexps []*regexp.Regexp
}

func (c *GinLoggerConfig) init() {
	if c.Formatter == nil {
		c.Formatter = defaultGinLogFormatter
	}
	if c.TraceIDFunc == nil {
		c.TraceIDFunc = defaultGinTraceIDFunc
	}

	var skip map[string]struct{}
	if length := len(c.SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, skipPath := range c.SkipPaths {
			skip[skipPath] = struct{}{}
		}
	}
	var skipRegexps []*regexp.Regexp
	for _, p := range c.SkipPathRegexps {
		if r, err := regexp.Compile(p); err != nil {
			Error("skip path regexps compile", zap.String("path_regex", p), zap.Error(err))
		} else {
			skipRegexps = append(skipRegexps, r)
		}
	}
	c.skipPathMap = skip
	c.skipRegexps = skipRegexps
	if c.SlowThreshold.Seconds() <= 0 {
		c.SlowThreshold = defaultGinSlowThreshold
	}
}

// GinLogger 以默认配置生成 gin 的 Logger 中间件
func GinLogger() gin.HandlerFunc {
	return GinLoggerWithConfig(GinLoggerConfig{})
}

// access log format for msg field
func defaultGinLogFormatter(m GinLogDetails) string {
	_, shortHandlerName := path.Split(m.HandlerName)
	msg := fmt.Sprintf("%s [%s] %s%s %d %f %s",
		m.ClientIP,
		m.Method,
		m.Host,
		m.RequestURI,
		m.StatusCode,
		m.Latency,
		shortHandlerName,
	)
	return msg
}

func defaultGinTraceIDFunc(c context.Context) (traceID string) {
	if c == nil {
		c = context.Background()
	}

	if gc, ok := c.(*gin.Context); ok {
		traceID = GetGinTraceIDFromHeader(gc)
		if traceID != "" {
			return
		}
		traceID = GetGinTraceIDFromQueryString(gc)
		if traceID != "" {
			return
		}
	}
	traceID = CtxTraceID(c)
	return
}

// GinLoggerWithConfig logger middleware for gin
func GinLoggerWithConfig(conf GinLoggerConfig) gin.HandlerFunc {
	conf.init()
	return func(c *gin.Context) {
		traceID := conf.TraceIDFunc(c)
		// set trace id to request header
		c.Request.Header.Set(TraceIDHeaderName, traceID)
		// set trace id to response header
		c.Writer.Header().Set(TraceIDHeaderName, traceID)
		// set trace id and ctxLogger to context
		ginLogger := CloneLogger("gin")
		if conf.InitFieldsFunc != nil {
			for k, v := range conf.InitFieldsFunc(c) {
				ginLogger = ginLogger.With(zap.Any(k, v))
			}
		}
		_, ctxLogger := NewCtxLogger(c, ginLogger, traceID)

		start := time.Now()

		// get request details
		details := GinLogDetails{
			Method:        c.Request.Method,
			Path:          c.Request.URL.Path,
			Query:         c.Request.URL.RawQuery,
			Proto:         c.Request.Proto,
			ContentLength: int(c.Request.ContentLength),
			Host:          c.Request.Host,
			RemoteAddr:    c.Request.RemoteAddr,
			RequestURI:    c.Request.RequestURI,
			Referer:       c.Request.Referer(),
			UserAgent:     c.Request.UserAgent(),

			ClientIP:    c.ClientIP(),
			ContentType: c.ContentType(),
			HandlerName: c.HandlerName(),
		}
		defer func() {
			// get response info
			details.StatusCode = c.Writer.Status()
			details.BodySize = c.Writer.Size()
			details.Timestamp = time.Now()
			details.Latency = details.Timestamp.Sub(start).Seconds()

			makeLog(ctxLogger, details, c, conf)
		}()

		c.Next()
	}
}

func makeLog(ctxLogger *zap.Logger, details GinLogDetails, c *gin.Context, conf GinLoggerConfig) {
	accessLogger := ctxLogger.Named("gin").With(
		zap.String("client_ip", details.ClientIP),
		zap.String("method", details.Method),
		zap.String("path", details.Path),
		zap.String("host", details.Host),
		zap.Int("status_code", details.StatusCode),
		zap.Float64("latency", details.Latency),
	)
	hasError := false
	if len(c.Errors) > 0 {
		hasError = true
		accessLogger = accessLogger.With(zap.String("context_errors", c.Errors.String()))
	}

	// details logger can print more details
	detailsLogger := accessLogger.Named("details").With(
		zap.String("query", details.Query),
		zap.String("proto", details.Proto),
		zap.Int("content_length", details.ContentLength),
		zap.String("remote_addr", details.RemoteAddr),
		zap.String("request_uri", details.RequestURI),
		zap.String("referer", details.Referer),
		zap.String("user_agent", details.UserAgent),
		zap.String("content_type", details.ContentType),
		zap.Int("body_size", details.BodySize),
		zap.String("handler_name", details.HandlerName),
	)

	logger := accessLogger
	// whether print details
	if conf.EnableDetails {
		logger = detailsLogger
	}

	// print access log
	log := logger.Info
	if details.StatusCode >= http.StatusInternalServerError || hasError {
		errLogger := detailsLogger.Named("err").With(
			zap.Any("context_keys", c.Keys),
			zap.Any("request_header", c.Request.Header),
			zap.Any("request_form", c.Request.Form),
			zap.String("request_body", string(GetGinRequestBody(c))),
		)
		if len(c.Errors) > 0 {
			for i, err := range c.Errors {
				errLogger = errLogger.With(zap.NamedError(fmt.Sprintf("error%d", i), err.Err))
			}
		}
		log = errLogger.Error
	} else if details.StatusCode >= http.StatusBadRequest {
		log = logger.Warn
		if len(c.Errors) > 0 {
			log = logger.Error
		}
	} else if len(c.Errors) > 0 {
		log = logger.Error
	}

	skipLog := false
	if _, exists := conf.skipPathMap[details.Path]; exists {
		skipLog = true
	} else {
		for _, p := range conf.skipRegexps {
			if p.MatchString(details.Path) {
				skipLog = true
				break
			}
		}
	}
	if !skipLog {
		// Warn log for slow request
		if details.Latency > conf.SlowThreshold.Seconds() {
			logger.Warn(
				conf.Formatter(details)+" hit slow request.",
				zap.Float64("slow_threshold", conf.SlowThreshold.Seconds()),
			)
		} else {
			log(conf.Formatter(details))
		}
	}
}

// GetGinRequestBody get request body
func GetGinRequestBody(c *gin.Context) []byte {
	var requestBody []byte
	if c.Request.Body != nil {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
		} else {
			requestBody = body
			// body is set to nil after read or bind, set it again
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
	return requestBody
}
