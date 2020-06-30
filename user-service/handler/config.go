package handler

import "github.com/x-punch/go-config"

// Config represents handler config
type Config struct {
	SiteURL         string          `toml:"site_url"`
	TokenSecret     string          `toml:"token_secret"`
	TokenExpiration config.Duration `toml:"token_expiration"`
}
