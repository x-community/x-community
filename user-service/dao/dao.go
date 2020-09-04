package dao

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/x-community/user-service/models"
	"gorm.io/gorm"
)

var _ UserDao = &userDao{}

type userDao struct {
	db *gorm.DB
}

func (d *userDao) IsEntityNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (d *userDao) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *userDao) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *userDao) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var entity models.User
	if err := d.db.WithContext(ctx).Find(&entity, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *userDao) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var entity models.User
	if err := d.db.Find(&entity, "username = ?", username).Error; err != nil {
		if d.IsEntityNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *userDao) EncryptPassword(ctx context.Context, password, salt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(salt+password+salt)))
}

func (d *userDao) CreateUser(ctx context.Context, tx *gorm.DB, entity *models.User) error {
	return tx.Create(entity).Error
}

func (d *userDao) ActiveUser(ctx context.Context, tx *gorm.DB, code string) error {
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
}

func (d *userDao) FellowUser(ctx context.Context, userID, fellowUserID uint32) error {
	if userID == fellowUserID {
		return nil
	}
	sql := "INSERT INTO `user_relation` (`user_id`, `fellow_user_id`, `created_at`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE created_at = VALUES(created_at);"
	if err := d.db.Exec(sql, userID, fellowUserID, time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (d *userDao) GetFellowUserCount(ctx context.Context, userID uint32) (uint32, error) {
	var count int64
	if err := d.db.Model(models.UserRelation{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return uint32(count), nil
}

func (d *userDao) GetFellowerCount(ctx context.Context, userID uint32) (uint32, error) {
	var count int64
	if err := d.db.Model(models.UserRelation{}).Where("fellow_user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return uint32(count), nil
}

func (d *userDao) GetFellowUsers(ctx context.Context, userID uint32, skip uint32, limit uint32) (uint32, []models.User, error) {
	total, err := d.GetFellowUserCount(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	var fellowUsers []models.User
	sql := "SELECT u.id, u.username, r.created_at FROM `user` as u JOIN `user_relation` as r ON u.id = r.fellow_user_id WHERE r.user_id = ?;"
	if err := d.db.Exec(sql, userID).Scan(&fellowUsers).Error; err != nil {
		return 0, nil, err
	}
	return total, nil, nil
}

func (d *userDao) GetFellowers(ctx context.Context, userID uint32, skip uint32, limit uint32) (uint32, []models.User, error) {
	total, err := d.GetFellowerCount(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	var fellowUsers []models.User
	sql := "SELECT u.id, u.username, r.created_at FROM `user` as u JOIN `user_relation` as r ON u.id = r.fellow_user_id WHERE u.id = ?;"
	if err := d.db.Joins(sql, userID).Scan(&fellowUsers).Error; err != nil {
		return 0, nil, err
	}
	return total, nil, nil
}

func (d *userDao) Transaction(fn func(*gorm.DB) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error { return fn(tx) })
}
