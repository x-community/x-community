package dao

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-community/user-service/models"
)

var _ UserDao = &userDao{}

type userDao struct {
	db *gorm.DB
}

func (d *userDao) IsEmailExists(email string) (bool, error) {
	var count int
	if err := d.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *userDao) IsUsernameExists(username string) (bool, error) {
	var count int
	if err := d.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *userDao) FindUserByEmail(email string) (*models.User, error) {
	var entity models.User
	if err := d.db.Find(&entity, "email = ?", email).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *userDao) FindUserByUsername(username string) (*models.User, error) {
	var entity models.User
	if err := d.db.Find(&entity, "username = ?", username).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *userDao) EncryptPassword(password, salt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(salt+password+salt)))
}

func (d *userDao) CreateUser(entity *models.User) error {
	return d.db.Create(entity).Error
}

func (d *userDao) ActiveUser(code string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		var entity models.User
		if err := tx.Where("active_code = ?", code).Find(&entity).Error; err != nil {
			return err
		}
		if !entity.Actived {
			updates := map[string]interface{}{"actived": true, "actived_at": time.Now().UTC()}
			if err := tx.Model(entity).Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *userDao) IsEntityNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
