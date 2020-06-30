package config

import (
	"os"

	"github.com/x-community/api-service/controller"
	"github.com/x-punch/go-config"
)

// Config represents the server configuration.
type Config struct {
	Address  string            `toml:"address"`
	Name     string            `toml:"-"`
	Version  string            `toml:"-"`
	Tracing  TracingConfig     `toml:"tracing"`
	API      controller.Config `toml:"api"`
	Services ServicesConfig    `toml:"services"`
}

// TracingConfig represents opentracing config
type TracingConfig struct {
	Enable    bool   `toml:"enable"`
	Collector string `toml:"collector"`
}

// ServicesConfig represents service config
type ServicesConfig struct {
	UserService string `toml:"user"`
}

// Load parse config info from config file and env args
func Load() (cfg *Config, err error) {
	cfg = &Config{
		Address: ":80",
		Name:    "api-service",
		Version: "0.0.0",
		API:     controller.NewConfig(),
	}
	if _, err := os.Stat("config.toml"); err == nil {
		err = config.Load(cfg, "config.toml")
	} else {
		err = config.Load(cfg)
	}
	return
}
