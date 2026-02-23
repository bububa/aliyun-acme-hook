package app

import (
	"github.com/urfave/cli/v2"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/service"
)

func beforeAction(c *cli.Context) error {
	var cfg config.Config
	cfgPath := c.String("config")
	if err := config.Load(cfgPath, &cfg); err != nil {
		return err
	}
	service.Init(c.Context, &cfg)
	return nil
}

func afterAction(c *cli.Context) error {
	return nil
}
