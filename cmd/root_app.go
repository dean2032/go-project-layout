package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              filepath.Base(os.Args[0]),
	Short:            "Clean architecture using gin framework",
	Long:             `This is a command runner or web server for api architecture in golang. `,
	TraverseChildren: true,
}

var (
	verbose = false
)

func SetAppVersion(version, commit, buildTime string) {
	rootCmd.Version = fmt.Sprintf("%s(rev:%s build time: %s)", version, commit, buildTime)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

// App root of application
type App struct {
	*cobra.Command
}

func NewApp() App {
	app := App{
		Command: rootCmd,
	}
	app.AddCommand(GetSubCommands(CommonModules)...)
	return app
}

var RootApp = NewApp()
