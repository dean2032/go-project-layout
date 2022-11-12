package controllers

import (
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
)

// EchoController struct
type EchoController struct {
}

// NesEchoController creates new controller
func NewEchoController() *EchoController {
	return &EchoController{}
}

// Echo echo user's input param
func (s *EchoController) Echo(c *gin.Context) {
	logging.Info("Echo route called")
	req := struct {
		Input string `form:"input" binding:"required"`
	}{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		OnError(c, err)
	}
	OnSuccess(c, req.Input)
}
