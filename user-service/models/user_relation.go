package models

import "time"

// UserRelation represents user relation
type UserRelation struct {
	ID           int64 `gorm:"primary_key"`
	UserID       uint32
	FellowUserID uint32    `gorm:"index:idx_fellow_uid;"`
	CreatedAt    time.Time `gorm:"index:idx_created_at;"`
}

// TableName represents database table name
func (UserRelation) TableName() string {
	return "user_relation"
}
