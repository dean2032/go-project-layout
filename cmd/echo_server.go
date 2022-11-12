package cmd

import (
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/api/routes"
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/spf13/cobra"
)

// EchoServerCommand test echo server command
type EchoServerCommand struct{}

func (s *EchoServerCommand) Short() string {
	return "simple echo server application"
}

func (s *EchoServerCommand) Setup(cmd *cobra.Command) {}

func (s *EchoServerCommand) Run() utils.CommandRunner {
	return func(
		cfg *config.Config,
		router *middlewares.RequestHandler,
		route *routes.EchoRoutes,
	) {
		route.Setup()

		logging.Infof("Running simple echo server on %s", cfg.ServerPort)
		if cfg.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + cfg.ServerPort)
		}
	}
}

func NewEchoServerCommand() *EchoServerCommand {
	return &EchoServerCommand{}
}
