package handler

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/server"
	"github.com/x-community/x-community/user-service/dao"
	pb "github.com/x-community/x-community/user-service/proto"
)

// Options represents handler options
type Options struct {
	ServiceName     string
	DB              *gorm.DB
	TokenSecret     string
	TokenExpiration time.Duration
	// MailService       pb.MailService
}

// RegisterHandlers will register handlers for user services
func RegisterHandlers(s server.Server, o Options) error {
	userDao := dao.NewUserDao(o.DB)
	err := pb.RegisterUserServiceHandler(s, &userService{id: o.ServiceName, dao: userDao})
	if err != nil {
		return err
	}
	return nil
}
