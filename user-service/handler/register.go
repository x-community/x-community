package handler

import (
	"gorm.io/gorm"
	"github.com/micro/go-micro/v2/server"
	"github.com/x-community/user-service/dao"
	pb "github.com/x-community/user-service/proto"
)

// Options represents handler options
type Options struct {
	ServiceName string
	Config      Config
	DB          *gorm.DB
	MailService pb.MailService
}

// Register will register handlers for user services
func Register(s server.Server, opts Options) error {
	userDao := dao.NewUserDao(opts.DB)
	userService := &userService{id: opts.ServiceName, cfg: opts.Config, dao: userDao, mailService: opts.MailService}
	err := pb.RegisterUserServiceHandler(s, userService)
	if err != nil {
		return err
	}
	err = pb.RegisterUserRelationServiceHandler(s, userService)
	if err != nil {
		return err
	}
	return nil
}
