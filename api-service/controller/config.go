package controller

// Config represents api config
type Config struct {
	Swagger SwaggerConfig `toml:"swagger"`
}

// SwaggerConfig represents swagger config
type SwaggerConfig struct {
	Enable bool   `toml:"enable"`
	Base   string `toml:"base"`
	Host   string `toml:"host"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Swagger: SwaggerConfig{Enable: false, Base: "/", Host: "localhost"},
	}
}
