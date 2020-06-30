package handler

// Config represents the server configuration.
type Config struct {
	Mail MailConfig `toml:"mail"`
}

// MailConfig represents mail server configuration
type MailConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Sender   string `toml:"sender"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{Mail: MailConfig{
		Sender: "X Community",
	}}
}
