package cmd

import (
	"fmt"

	"github.com/neploxaudit/pocs-cli/internal/build"
	"github.com/urfave/cli/v2"
)

var Version = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "Shows version information about the CLI",
	Action: func(c *cli.Context) error {
		fmt.Printf("pocs cli version %s\n", build.Version)
		fmt.Printf("build commit %s\n", build.Commit)
		fmt.Printf("build time %s\n", build.Date)
		return nil
	},
}
