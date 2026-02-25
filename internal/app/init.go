package app

import (
	"github.com/urfave/cli/v2"
)

func NewApp(app *cli.App) {
	*app = cli.App{
		Name:    "aliyun-acme-hook",
		Version: "v1.2.0",
		Usage:   "aliyun acme deploy hook",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Load configuration from `FILE`",
				Required: true,
			},
		},
		Before: beforeAction,
		After:  afterAction,
		Commands: []*cli.Command{
			{
				Name:     "certificate",
				Usage:    "certificate update in aliyun",
				Category: "Aliyun",
				Action:   Deploy,
			},
		},
	}
}
