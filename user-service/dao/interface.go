package dao

import (
	"context"

	"github.com/x-community/user-service/models"
	"gorm.io/gorm"
)

// UserDao represents user data access service
type UserDao interface {
	IsEntityNotFoundError(err error) bool
	Transaction(fn func(*gorm.DB) error) error

	IsEmailExists(ctx context.Context, email string) (bool, error)
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	EncryptPassword(ctx context.Context, password, salt string) string
	CreateUser(ctx context.Context, tx *gorm.DB, user *models.User) error
	ActiveUser(ctx context.Context, tx *gorm.DB, code string) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)

	FellowUser(ctx context.Context, userID uint32, fellowerID uint32) error
	GetFellowUserCount(ctx context.Context, userID uint32) (uint32, error)
	GetFellowerCount(ctx context.Context, userID uint32) (uint32, error)
	GetFellowUsers(ctx context.Context, userID uint32, skip uint32, limit uint32) (uint32, []models.User, error)
	GetFellowers(ctx context.Context, userID uint32, skip uint32, limit uint32) (uint32, []models.User, error)
}

// NewUserDao will create user dao instance
func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{db: db}
}
