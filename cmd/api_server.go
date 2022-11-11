package cmd

import (
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/api/routes"
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/repo"
	"github.com/dean2032/go-project-layout/utils"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/spf13/cobra"
)

// ApiServerCommand test command
type ApiServerCommand struct{}

func (s *ApiServerCommand) Short() string {
	return "serve application"
}

func (s *ApiServerCommand) Setup(cmd *cobra.Command) {}

func (s *ApiServerCommand) Run() utils.CommandRunner {
	return func(
		middleware middlewares.Middlewares,
		cfg *config.Config,
		router *middlewares.RequestHandler,
		route routes.ApiRoutes,
		database *repo.Database,
	) {
		middleware.Setup()
		route.Setup()

		logging.Info("Running api server")
		if cfg.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + cfg.ServerPort)
		}
	}
}

func NewApiServerCommand() *ApiServerCommand {
	return &ApiServerCommand{}
}
