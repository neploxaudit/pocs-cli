package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/neploxaudit/pocs-cli/internal/cmd"
	"github.com/neploxaudit/pocs-cli/internal/config"
	"github.com/neploxaudit/pocs-cli/pkg/pocs"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "pocs",
	Usage: "Convenient CLI for pocs.neplox.security",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "api",
			Usage: "Base `URL` for the pocs API",
			Value: "https://pocs.neplox.security",
		},
	},
	Before: func(ctx *cli.Context) error {
		baseURL, err := url.Parse(ctx.String("api"))
		if err == nil && baseURL.Scheme != "http" && baseURL.Scheme != "https" {
			err = fmt.Errorf("scheme must be http or https")
		}
		if err != nil {
			return &cmd.UsageError{Message: fmt.Sprintf("--api value is not a valid URL: %s", err)}
		}

		config.NeploxToken = os.Getenv("NEPLOX_TOKEN")
		config.BaseURL = baseURL
		config.Client = pocs.NewClient(config.BaseURL, config.NeploxToken)

		return nil
	},
	Commands: []*cli.Command{
		cmd.Version,
		cmd.Read,
	},
	HideVersion: true,
	CommandNotFound: func(ctx *cli.Context, s string) {
		fmt.Printf("No command found for %s\n\n", s)
		cli.ShowAppHelpAndExit(ctx, 1)
	},
}

func main() {
	mainCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	appCtx := cli.NewContext(app, nil, &cli.Context{Context: mainCtx})
	runErr := app.RunContext(mainCtx, os.Args)

	if uerr, ok := lo.ErrorsAs[*cmd.UsageError](runErr); ok {
		fmt.Printf("Incorrect usage: %s\n\n", uerr.Message)

		if uerr.Command == "" {
			cli.ShowAppHelpAndExit(appCtx, 1)
		} else {
			cli.ShowCommandHelpAndExit(appCtx, uerr.Command, 1)
		}
	} else if eerr, ok := lo.ErrorsAs[*cmd.ExecError](runErr); ok {
		fmt.Printf("Failed to `%s`: %s\n\n", eerr.Command, eerr.Message)
		os.Exit(1)
	} else if runErr != nil {
		os.Exit(1)
	}
}
