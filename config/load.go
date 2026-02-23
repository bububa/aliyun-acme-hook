package config

import (
	"github.com/jinzhu/configor"
)

func Load(configPath string, cfg *Config) error {
	loader := configor.New(&configor.Config{
		Verbose:              false,
		ErrorOnUnmatchedKeys: true,
		Environment:          "production",
	})
	return loader.Load(cfg, configPath)
}
