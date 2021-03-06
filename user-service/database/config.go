package database

// Config represents the mysql configuration.
type Config struct {
	Host               string `toml:"host"`
	Port               int    `toml:"port"`
	Database           string `toml:"database"`
	Charset            string `toml:"charset"`
	Username           string `toml:"username"`
	Password           string `toml:"password"`
	MaxIdleConnections int    `toml:"max_idle"`
	MaxOpenConnections int    `toml:"max_open"`
	ShowSQL            bool   `toml:"show_sql"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Host:               "127.0.0.1",
		Port:               3306,
		Charset:            "utf8",
		MaxIdleConnections: 5,
		MaxOpenConnections: 10,
		ShowSQL:            false,
	}
}
