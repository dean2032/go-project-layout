package routes

import (
	"github.com/dean2032/go-project-layout/api/controllers"
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/utils/logging"
)

// UserRoutes struct
type UserRoutes struct {
	handler        *middlewares.RequestHandler
	userController *controllers.UserController
	authMiddleware *middlewares.JWTAuthMiddleware
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	handler *middlewares.RequestHandler,
	userController *controllers.UserController,
	authMiddleware *middlewares.JWTAuthMiddleware,
) *UserRoutes {
	return &UserRoutes{
		handler:        handler,
		userController: userController,
		authMiddleware: authMiddleware,
	}
}

// Setup user routes
func (s *UserRoutes) Setup() {
	logging.Info("Setting up routes")
	api := s.handler.Gin.Group("/api").Use(s.authMiddleware.Handler())
	{
		api.GET("/user", s.userController.GetUser)
		api.GET("/user/:id", s.userController.GetOneUser)
		api.POST("/user", s.userController.SaveUser)
		api.POST("/user/:id", s.userController.UpdateUser)
		api.DELETE("/user/:id", s.userController.DeleteUser)
	}
}
