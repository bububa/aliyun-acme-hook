package service

import (
	"github.com/bububa/aliyun-acme-hook/config"
)

var configStore *config.Config

func Config() *config.Config {
	return configStore
}
