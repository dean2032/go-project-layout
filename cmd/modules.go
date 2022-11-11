package cmd

import (
	"github.com/dean2032/go-project-layout/api/controllers"
	"github.com/dean2032/go-project-layout/api/middlewares"
	"github.com/dean2032/go-project-layout/api/routes"
	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/repo"
	"github.com/dean2032/go-project-layout/services"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	controllers.Module,
	routes.Module,
	services.Module,
	middlewares.Module,
	repo.Module,
	config.Module,
)
