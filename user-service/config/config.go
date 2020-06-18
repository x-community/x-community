package config

import (
	"os"
	"time"

	"github.com/x-community/x-community/user-service/database"
	"github.com/x-punch/go-config"
)

// Config represents the server configuration.
type Config struct {
	Address  string          `toml:"address"`
	Name     string          `toml:"name"`
	Version  string          `toml:"-"`
	DB       database.Config `toml:"db"`
	Tracing  TracingConfig   `toml:"tracing"`
	Services Services        `toml:"services"`
}

// TracingConfig represents opentracing config
type TracingConfig struct {
	Enable    bool   `toml:"enable"`
	Collector string `toml:"collector"`
}

// Services represents third-services
type Services struct {
	TokenSecret     string          `toml:"token_secret"`
	TokenExpiration config.Duration `toml:"token_expiration"`
	MailService     string          `toml:"mail_service"`
}

// Load parse config info from config file and env args
func Load() (cfg *Config, err error) {
	cfg = &Config{
		Address: ":80",
		Name:    "user-service",
		Version: "1.0.0",
		Services: Services{
			TokenExpiration: config.Duration(time.Hour),
			TokenSecret:     "secret",
		},
		DB: database.NewConfig(),
	}
	if _, err := os.Stat("config.toml"); err == nil {
		err = config.Load(cfg, "config.toml")
	} else {
		err = config.Load(cfg)
	}
	return
}
