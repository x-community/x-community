package handler

import (
	"context"
	"regexp"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/x-community/x-community/user-service/dao"
	"github.com/x-community/x-community/user-service/models"
	pb "github.com/x-community/x-community/user-service/proto"
	"github.com/x-punch/go-strings"
)

type userService struct {
	id  string
	dao dao.UserDao
}

func (s *userService) Register(ctx context.Context, in *pb.RegisterRequest, out *pb.RegisterResponse) error {
	if len(in.Email) == 0 || len(in.Email) > 256 || !regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`).MatchString(in.Email) {
		return s.NewError(errInvalidEmail)
	}
	if len(in.Username) == 0 || len(in.Username) > 32 {
		return s.NewError(errInvalidUsername)
	}
	if len(in.Password) < 8 || !regexp.MustCompile("^([A-Z]|[a-z]|[0-9]|[`~!@#$%^&*()-_+=|{}':;\\\\[\\]<>,./?]){8,}$").MatchString(in.Password) {
		return s.NewError(errInvalidPassword)
	}
	if exists, err := s.dao.IsEmailExists(in.Email); err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	} else if exists {
		return s.NewError(errEmailAlreadyRegistered)
	}
	if exists, err := s.dao.IsUsernameExists(in.Email); err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	} else if exists {
		return s.NewError(errUsernameAlreadyRegistered)
	}
	salt := strings.GetRandomString(10)
	user := &models.User{
		Email:    in.Email,
		Username: in.Username,
		Salt:     salt,
		Password: s.dao.EncryptPassword(in.Password, salt),
	}
	if err := s.dao.CreateUser(user); err != nil {
		return s.InternalServerError(err.Error())
	}
	return nil
}

func (s *userService) Authenticate(ctx context.Context, in *pb.AuthenticateRequest, out *pb.AuthenticateResponse) error {
	if len(in.EmailOrUsername) == 0 || len(in.Password) == 0 {
		return s.NewError(errIncorrectUsernameOrPassword)
	}
	user, err := s.dao.FindUserByEmail(in.EmailOrUsername)
	if err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	}
	if user == nil {
		user, err = s.dao.FindUserByUsername(in.EmailOrUsername)
		if err != nil {
			log.Error(err)
			return s.InternalServerError(err.Error())
		}
		if user == nil {
			return s.NewError(errIncorrectUsernameOrPassword)
		}
	}
	if user.Password != s.dao.EncryptPassword(in.Password, user.Salt) {
		return s.NewError(errIncorrectUsernameOrPassword)
	}
	return nil
}
