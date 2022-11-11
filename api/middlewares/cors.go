package middlewares

import (
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils/logging"
	cors "github.com/rs/cors/wrapper/gin"
)

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	handler *RequestHandler
	cfg     *config.Config
}

// NewCorsMiddleware creates new cors middleware
func NewCorsMiddleware(handler *RequestHandler, cfg *config.Config) *CorsMiddleware {
	return &CorsMiddleware{
		handler: handler,
		cfg:     cfg,
	}
}

// Setup sets up cors middleware
func (m *CorsMiddleware) Setup() {
	logging.Info("Setting up cors middleware")

	m.handler.Gin.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		Debug:            m.cfg.Debug,
	}))
}
