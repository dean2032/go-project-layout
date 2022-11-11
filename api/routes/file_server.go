package routes

import (
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
)

// FileRoutes struct
type FileRoutes struct {
	handler *middlewares.RequestHandler
	cfg     *config.Config
}

// NewUserRoutes creates new user controller
func NewFileRoutes(
	handler *middlewares.RequestHandler,
	cfg *config.Config,
) *FileRoutes {
	return &FileRoutes{
		handler: handler,
		cfg:     cfg,
	}
}

// Setup user routes
func (s *FileRoutes) Setup() {
	logging.Infof("Setting up file routes on cfg.PublicDir: %s", s.cfg.PublicDir)
	s.handler.Gin.StaticFS("/", gin.Dir(s.cfg.PublicDir, true))
}
