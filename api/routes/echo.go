package routes

import (
	"github.com/dean2032/go-project-layout/api/controllers"
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils/logging"
)

// EchoRoutes struct
type EchoRoutes struct {
	handler *middlewares.RequestHandler
	cfg     *config.Config
	echo    *controllers.EchoController
}

// NewUserRoutes creates new user controller
func NewEchoRoutes(
	handler *middlewares.RequestHandler,
	cfg *config.Config,
	echo *controllers.EchoController,
) *EchoRoutes {
	return &EchoRoutes{
		handler: handler,
		cfg:     cfg,
		echo:    echo,
	}
}

// Setup user routes
func (s *EchoRoutes) Setup() {
	logging.Infof("Setting up echo routes")
	s.handler.Gin.GET("/echo", s.echo.Echo)
}
