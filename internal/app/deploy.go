package app

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v2"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/cas"
	"github.com/bububa/aliyun-acme-hook/internal/cdn"
	"github.com/bububa/aliyun-acme-hook/internal/model"
	"github.com/bububa/aliyun-acme-hook/internal/service"
	"github.com/bububa/aliyun-acme-hook/internal/slb"
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
	cert := new(model.Cert)
	if err := cert.LoadFromEnv(); err != nil {
		slog.ErrorContext(ctx, "load certification info from env failed", "error", err)
		return err
	}
	if err := cas.Upload(ctx, cfg.CAS, cert); err != nil {
		slog.ErrorContext(ctx, "upload to CAS failed", "account", cfg.Name, "error", err)
		return err
	}
	slog.InfoContext(ctx, "uploaded to CAS", "account", cfg.Name, "certID", cert.ID)
	if cfg.CDN != nil {
		if err := cdn.Certificate(ctx, cfg.CDN, cert); err != nil {
			slog.ErrorContext(ctx, "update CDN certification failed", "account", cfg.Name, "error", err)
			return err
		} else {
			slog.InfoContext(ctx, "updated CDN certification", "account", cfg.Name)
			return nil
		}
	}
	if cfg.SLB != nil {
		if err := slb.Certificate(ctx, cfg.SLB, cert); err != nil {
			slog.ErrorContext(ctx, "update SLB certification failed", "account", cfg.Name, "error", err)
			return err
		} else {
			slog.InfoContext(ctx, "updated SLB certification", "account", cfg.Name)
			return nil
		}
	}
	return nil
}
