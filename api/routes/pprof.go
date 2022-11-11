package routes

import (
	"net/http/pprof"

	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
)

// PprofRoutes struct
type PprofRoutes struct {
	handler *middlewares.RequestHandler
	cfg     *config.Config
}

// NewUserRoutes creates new user controller
func NewPprofRoutes(
	handler *middlewares.RequestHandler,
	cfg *config.Config,
) *PprofRoutes {
	return &PprofRoutes{
		handler: handler,
		cfg:     cfg,
	}
}

// Setup user routes
func (s *PprofRoutes) Setup() {
	if s.cfg.PprofPath != "" {
		logging.Infof("Setting up pprof routes on %s", s.cfg.PprofPath)
		registerPprof(s.handler.Gin.Group(s.cfg.PprofPath))
	}
}

func registerPprof(router *gin.RouterGroup) {
	router.GET("/", gin.WrapF(pprof.Index))
	router.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	router.GET("/profile", gin.WrapF(pprof.Profile))
	router.POST("/symbol", gin.WrapF(pprof.Symbol))
	router.GET("/symbol", gin.WrapF(pprof.Symbol))
	router.GET("/trace", gin.WrapF(pprof.Trace))
	router.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	router.GET("/block", gin.WrapH(pprof.Handler("block")))
	router.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	router.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	router.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	router.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
}
