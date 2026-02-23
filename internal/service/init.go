package service

import (
	"context"

	"github.com/bububa/aliyun-acme-hook/config"
)

func Init(ctx context.Context, cfg *config.Config) {
	configStore = cfg
}
