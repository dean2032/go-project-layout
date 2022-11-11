package routes

import (
	"github.com/dean2032/go-project-layout/api/controllers"
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/utils/logging"
)

// AuthRoutes struct
type AuthRoutes struct {
	handler        *middlewares.RequestHandler
	authController *controllers.JWTAuthController
}

// NewAuthRoutes creates new user controller
func NewAuthRoutes(
	handler *middlewares.RequestHandler,
	authController *controllers.JWTAuthController,
) *AuthRoutes {
	return &AuthRoutes{
		handler:        handler,
		authController: authController,
	}
}

// Setup user routes
func (s *AuthRoutes) Setup() {
	logging.Info("Setting up routes")
	auth := s.handler.Gin.Group("/auth")
	{
		auth.POST("/login", s.authController.SignIn)
		auth.POST("/register", s.authController.Register)
	}
}
