package models

import "time"

type User struct {
	ID               uint32     `gorm:"primary_key"`
	Username         string     `gorm:"type:varchar(32);not null;"`
	Email            string     `gorm:"type:varchar(256);unique_index;"`
	IsEmailConfirmed bool       `gorm:"type:tinyint(1)"`
	Salt             string     `gorm:"type:varchar(10);"`
	Password         string     `gorm:"type:varchar(128);"`
	LastLoginIP      string     `gorm:"type:varchar(64);"`
	LastLoginTime    *time.Time `gorm:"type:datetime"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TableName represents database table name
func (User) TableName() string {
	return "users"
}
