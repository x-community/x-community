package models

import "time"

// User represents user
type User struct {
	ID            uint32 `gorm:"primary_key"`
	Username      string `gorm:"type:varchar(32);unique_index"`
	Email         string `gorm:"type:varchar(256);unique_index;"`
	Salt          string `gorm:"type:varchar(10);"`
	Password      string `gorm:"type:varchar(128);"`
	Actived       bool
	ActiveCode    string     `gorm:"type:varchar(32);"`
	LastLoginIP   string     `gorm:"type:varchar(64);"`
	LastLoginTime *time.Time `gorm:"type:datetime"`
	CreatedAt     time.Time
	ActivedAt     *time.Time
	UpdatedAt     time.Time
}

// TableName represents database table name
func (User) TableName() string {
	return "user"
}
