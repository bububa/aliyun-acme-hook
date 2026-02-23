package app

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v2"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/cas"
	"github.com/bububa/aliyun-acme-hook/internal/service"
)

func Deploy(c *cli.Context) error {
	ctx := c.Context
	for _, cfg := range service.Config().Accounts {
		if err := AccountDeploy(ctx, &cfg); err != nil {
			// TODO
		}
	}
	return nil
}

func AccountDeploy(ctx context.Context, cfg *config.Account) error {
	if certID, err := cas.Upload(ctx, cfg.CAS); err != nil {
		slog.ErrorContext(ctx, "upload to cas failedw", "error", err)
		return err
	} else {
		slog.InfoContext(ctx, "uploaded to cas", "certID", certID)
		return nil
	}
}
