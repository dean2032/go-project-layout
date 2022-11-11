package main

import (
	"github.com/dean2032/go-project-layout/cmd"
	_ "github.com/go-sql-driver/mysql"
)

var (
	BuildTime    = "2022-11-11.11.11.11"
	GitCommit    = "na"
	BuildVersion = "v0.0.1"
)

func main() {
	cmd.SetAppVersion(BuildVersion, GitCommit, BuildTime)
	cmd.RootApp.Execute()
}
