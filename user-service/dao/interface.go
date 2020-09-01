package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/x-community/user-service/models"
)

// UserDao represents user data access service
type UserDao interface {
	IsEmailExists(email string) (bool, error)
	IsUsernameExists(username string) (bool, error)
	EncryptPassword(password, salt string) string
	CreateUser(*models.User) error
	ActiveUser(code string) error
	FindUserByEmail(string) (*models.User, error)
	FindUserByUsername(string) (*models.User, error)
	IsEntityNotFoundError(error) bool

	FellowUser(uint32, uint32) error
	GetFellowUserCount(uint32) (uint32, error)
	GetFellowerCount(uint32) (uint32, error)
	GetFellowUserIDs(uint32, uint32, uint32) ([]uint32, error)
	GetFellowerIDs(uint32, uint32, uint32) ([]uint32, error)
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{db: db}
}
