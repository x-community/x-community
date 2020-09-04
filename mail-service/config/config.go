package config

import (
	"os"

	"github.com/x-community/mail-service/handler"
	"github.com/x-punch/go-config"
)

// Config represents the server configuration.
type Config struct {
	Address string         `toml:"address"`
	Name    string         `toml:"-"`
	Version string         `toml:"-"`
	Tracing TracingConfig  `toml:"tracing"`
	Handler handler.Config `toml:"handler"`
}

// TracingConfig represents opentracing config
type TracingConfig struct {
	Enable    bool   `toml:"enable"`
	Collector string `toml:"collector"`
}

// Load parse config info from config file and env args
func Load() (cfg *Config, err error) {
	cfg = &Config{
		Address: ":80",
		Name:    "mail-service",
		Version: "0.0.1",
	}
	if _, err := os.Stat("config.toml"); err == nil {
		err = config.Load(cfg, "config.toml")
	} else {
		err = config.Load(cfg)
	}
	return
}
