package cmd

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/neploxaudit/pocs-cli/internal/config"
	"github.com/urfave/cli/v2"
)

var Read = &cli.Command{
	Name:      "read",
	Aliases:   []string{"r"},
	Usage:     "Read file or directory from pocs, by default displaying the contents with syntax highlighting",
	UsageText: "pocs read <path | link> [command options]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "style",
			Usage:   "Chroma (github.com/alecthomas/chroma) style name to use for rendering output",
			Value:   "evergarden",
			EnvVars: []string{"POCS_STYLE"},
		},
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			return &UsageError{Command: c.Command.Name, Message: "either path (pocsdir/.../poc.html) or link (https://pocs.neplox.security/pocsdir/.../poc.html) is required"}
		}

		// Parse
		link, err := url.Parse(c.Args().First())
		if err != nil {
			return &UsageError{Command: c.Command.Name, Message: fmt.Sprintf("invalid link: %s", err)}
		}

		if link.Scheme != "" || link.Host != "" {
			if link.Scheme != config.BaseURL.Scheme || link.Host != config.BaseURL.Host {
				return &UsageError{Command: c.Command.Name, Message: fmt.Sprintf("specified link %s is not for the configured base URL (%s)", link, config.BaseURL)}
			}
		}

		// Request
		rawBody, mime, err := config.Client.Get(link.Path)
		if err != nil {
			return &ExecError{Command: c.Command.Name, Message: fmt.Sprintf("fetching %s from API: %s", strconv.Quote(link.Path), err)}
		}

		body := strings.TrimSpace(string(rawBody))

		// Render
		lexer := lexers.MatchMimeType(mime)
		if lexer == nil {
			lexer = lexers.Match(link.Path)
		}
		if lexer == nil {
			lexer = lexers.Analyse(body)
		}
		if lexer == nil {
			lexer = lexers.Fallback
		}

		style := styles.Get(c.String("style"))
		if style == nil {
			style = styles.Fallback
		}

		formatter := formatters.Get("terminal256")
		if formatter == nil {
			formatter = formatters.Fallback
		}

		iterator, err := lexer.Tokenise(nil, body)
		if err != nil {
			fmt.Println(body)
			return nil
		}

		buf := bytes.NewBuffer(nil)
		if err := formatter.Format(buf, style, iterator); err != nil {
			fmt.Println(body)
			return nil
		}

		fmt.Println(buf.String())

		return nil
	},
}
