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

func (d *userDao) FellowUser(userID, fellowUserID uint32) error {
	if userID == fellowUserID {
		return nil
	}
	sql := "INSERT INTO `user_relation` (`user_id`, `fellow_user_id`, `created_at`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE created_at = VALUES(created_at);"
	if err := d.db.Exec(sql, userID, fellowUserID, time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (d *userDao) GetFellowUserCount(userID uint32) (uint32, error) {
	var count uint32
	if err := d.db.Model(models.UserRelation{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *userDao) GetFellowerCount(userID uint32) (uint32, error) {
	var count uint32
	if err := d.db.Model(models.UserRelation{}).Where("fellow_user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *userDao) GetFellowUsers(userID uint32, skip uint32, limit uint32) ([]uint32, error) {
	total, err := d.GetFellowUserCount(userID)
	if err != nil {
		return nil, err
	}
	type fellowUser struct {
		Username  string
		CreatedAt time.Time
	}
	fellowUsers
	sql := "SELECT u.username, r.created_at FROM `user` as u JOIN `user_relation` as r ON u.id = r.fellow_user_id WHERE u.id = ?;"
	if err := d.db.Exec(sql, userID).Scan().Error; err != nil {

	}
	return nil, nil
}

func (d *userDao) GetFellowerIDs(uint32, uint32, uint32) ([]uint32, error) {
	return nil, nil
}
