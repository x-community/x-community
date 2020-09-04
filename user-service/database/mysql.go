package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase return a new mysql db connection
func NewDatabase(cfg Config) (db *gorm.DB, err error) {
	gormConfig := &gorm.Config{
		Logger:  logger.Default.LogMode(cfg.logLevel()),
		NowFunc: func() time.Time { return time.Now().UTC() },
	}
	return gorm.Open(mysql.Open(cfg.dsn()), gormConfig)
}

func (c *Config) dsn() string {
	username := c.Username
	password := c.Password
	host := c.Host
	port := c.Port
	database := c.Database
	charset := c.Charset
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=UTC", username, password, host, port, database, charset)
}

func (c *Config) logLevel() logger.LogLevel {
	switch c.LogLevel {
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "silent":
		return logger.Silent
	default:
		return logger.Error
	}
}
