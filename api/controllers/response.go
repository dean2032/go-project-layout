package controllers

import (
	"net/http"

	"github.com/dean2032/go-project-layout/constants"
	"github.com/dean2032/go-project-layout/utils/errors"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
)

// Response common api response struct
type Response struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

// ApiHandler is a api handler function
type ApiHandler = func(c *gin.Context) *Response

// OnError make a error response
func OnError(c *gin.Context, err error) {
	cause := errors.Cause(err)
	response := &Response{Code: 100}
	var codeErr *errors.CodeError
	if realCodeErr, ok := cause.(*errors.CodeError); ok {
		codeErr = realCodeErr
	} else {
		codeErr = errors.Err2Code(err)
	}
	response.Code = codeErr.Code()
	response.Message = err.Error()
	c.Set(constants.ErrorCodeGinContextKey, codeErr.Error())
	setResponse(c, response)
}

// OnSuccess make a success response
func OnSuccess(c *gin.Context, data interface{}) {
	setResponse(c, &Response{
		Code:    0,
		Message: "OK",
		Data:    data,
	})
}

func setResponse(c *gin.Context, r *Response) {
	logger := logging.CtxLogger(c).Sugar()
	if !c.Writer.Written() {
		c.JSON(http.StatusOK, r)
	} else {
		logger.Warnf("get response but already has response body:%+v", r)
	}
}
