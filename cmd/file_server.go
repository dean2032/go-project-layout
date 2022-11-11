package cmd

import (
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/api/routes"
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/spf13/cobra"
)

// ServeCommand test command
type FileServerCommand struct{}

func (s *FileServerCommand) Short() string {
	return "simple file server application"
}

func (s *FileServerCommand) Setup(cmd *cobra.Command) {}

func (s *FileServerCommand) Run() utils.CommandRunner {
	return func(
		cfg *config.Config,
		router *middlewares.RequestHandler,
		route *routes.FileRoutes,
	) {
		route.Setup()

		logging.Infof("Running simple file server on %s", cfg.ServerPort)
		if cfg.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + cfg.ServerPort)
		}
	}
}

func NewFileServerCommand() *FileServerCommand {
	return &FileServerCommand{}
}
