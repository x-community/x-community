package database

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// NewDatabase return a new mysql db connection
func NewDatabase(config Config) (db *gorm.DB, err error) {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
	db, err = gorm.Open("mysql", connectionString(config))
	if nil == err && config.ShowSQL {
		db.LogMode(true)
	}
	return
}

func connectionString(config Config) string {
	username := config.Username
	password := config.Password
	host := config.Host
	port := config.Port
	database := config.Database
	charset := config.Charset
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=UTC", username, password, host, port, database, charset)
}
