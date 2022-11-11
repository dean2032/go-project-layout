package cmd

import (
	"context"
	"log"

	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var cmds = map[string]utils.Command{
	"file_server": NewFileServerCommand(),
	"api_serverr": NewApiServerCommand(),
}

// GetSubCommands gives a list of sub commands
func GetSubCommands(opt fx.Option) []*cobra.Command {
	var subCommands []*cobra.Command
	for name, cmd := range cmds {
		subCommands = append(subCommands, WrapSubCommand(name, cmd, opt))
	}
	return subCommands
}

func WrapSubCommand(name string, cmd utils.Command, opt fx.Option) *cobra.Command {
	wrappedCmd := &cobra.Command{
		Use:   name,
		Short: cmd.Short(),
		Run: func(c *cobra.Command, args []string) {
			others := fx.Options(
				fx.NopLogger,
				config.GenLoggerModule("access"),
			)
			if verbose {
				others = fx.Options()
			}
			opts := fx.Options(
				others,
				fx.Invoke(cmd.Run()),
			)
			ctx := context.Background()
			app := fx.New(opt, opts)
			err := app.Start(ctx)
			defer app.Stop(ctx)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Setup(wrappedCmd)
	return wrappedCmd
}
