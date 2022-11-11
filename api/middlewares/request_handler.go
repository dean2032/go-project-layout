package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/constants"
	"github.com/dean2032/go-project-layout/utils/errors"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(cfg *config.Config) *RequestHandler {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.New()
	app.Use(logging.GinLoggerWithConfig(logging.GinLoggerConfig{
		SkipPaths:     []string{"/"},
		EnableDetails: cfg.Debug,
		SlowThreshold: 10 * time.Second,
	}))
	app.Use(globalPanicHandler())
	app.NoMethod(handleNotFound)
	app.NoRoute(handleNotFound)
	return &RequestHandler{Gin: app}
}

func handleNotFound(c *gin.Context) {
	c.JSON(404, gin.H{"code": 404, "message": "not found"})
}

func globalPanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logging.CtxLogger(c).Error(fmt.Sprintf("panic: %v", err), zap.String("stack", string(stack)))
				// do this when response content not set
				if !c.Writer.Written() {
					response := gin.H{}
					if e, ok := err.(error); ok {
						response = parseError(c, e)
					} else {
						response = parseError(c, fmt.Errorf("%s", err))
					}
					c.JSON(http.StatusOK, response)
				}
				return
			}
		}()
		c.Next()
	}
}

func parseError(c *gin.Context, err error) gin.H {
	cause := errors.Cause(err)
	resp := gin.H{"code": 100, "message": "Unknown error"}
	var codeErr *errors.CodeError
	if realCodeErr, ok := cause.(*errors.CodeError); ok {
		codeErr = realCodeErr
	} else {
		codeErr = errors.Err2Code(err)
	}
	resp["code"] = codeErr.Code()
	resp["message"] = err.Error()
	c.Set(constants.ErrorCodeGinContextKey, codeErr.Error())
	return resp
}
