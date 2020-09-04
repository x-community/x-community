package handler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/matcornic/hermes/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/x-community/user-service/dao"
	"github.com/x-community/user-service/models"
	pb "github.com/x-community/user-service/proto"
	"github.com/x-punch/go-strings"
	"gorm.io/gorm"
)

const (
	// SaltLength represents salt max length in dataase
	SaltLength = 10
	// ActiveCodeLength represents active code max length in dataase
	ActiveCodeLength = 32
)

type userService struct {
	id          string
	cfg         Config
	dao         dao.UserDao
	mailService pb.MailService
}

func (s *userService) Register(ctx context.Context, in *pb.RegisterRequest, out *pb.RegisterReply) error {
	if len(in.Email) == 0 || len(in.Email) > 256 || !regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`).MatchString(in.Email) {
		return s.NewError(errInvalidEmail)
	}
	if len(in.Username) == 0 || len(in.Username) > 32 {
		return s.NewError(errInvalidUsername)
	}
	if len(in.Password) < 8 || !regexp.MustCompile("^([A-Z]|[a-z]|[0-9]|[`~!@#$%^&*()-_+=|{}':;\\\\[\\]<>,./?]){8,}$").MatchString(in.Password) {
		return s.NewError(errInvalidPassword)
	}
	if exists, err := s.dao.IsEmailExists(ctx, in.Email); err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	} else if exists {
		return s.NewError(errEmailAlreadyRegistered)
	}
	if exists, err := s.dao.IsUsernameExists(ctx, in.Username); err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	} else if exists {
		return s.NewError(errUsernameAlreadyRegistered)
	}
	salt := strings.GetRandomString(SaltLength)
	user := &models.User{
		Email:      in.Email,
		Username:   in.Username,
		Salt:       salt,
		Password:   s.dao.EncryptPassword(ctx, in.Password, salt),
		Actived:    false,
		ActiveCode: strings.GetRandomString(ActiveCodeLength),
	}
	if err := s.dao.Transaction(func(tx *gorm.DB) error {
		if err := s.dao.CreateUser(ctx, tx, user); err != nil {
			return err
		}
		return s.sendActivationEmail(user)
	}); err != nil {
		return s.InternalServerError(err.Error())
	}
	return nil
}

func (s *userService) Authenticate(ctx context.Context, in *pb.AuthenticateRequest, out *pb.AuthenticateReply) error {
	if len(in.Email) == 0 || len(in.Password) == 0 {
		return s.NewError(errIncorrectUsernameOrPassword)
	}
	user, err := s.dao.FindUserByEmail(ctx, in.Email)
	if err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	}
	if user.Password != s.dao.EncryptPassword(ctx, in.Password, user.Salt) {
		return s.NewError(errIncorrectUsernameOrPassword)
	}
	if !user.Actived {
		return s.NewError(errAccountNotActived)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        strconv.FormatUint(uint64(user.ID), 10),
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.TokenExpiration)).Unix(),
		Issuer:    "X Community",
	})
	out.Token, err = token.SignedString([]byte(s.cfg.TokenSecret))
	if err != nil {
		return s.InternalServerError(err.Error())
	}
	return nil
}

func (s *userService) VerifyAccount(ctx context.Context, in *pb.VerifyAccountRequest, out *pb.VerifyAccountReply) error {
	if len(in.Code) == 0 {
		return s.NewError(errInvalidActiveCode)
	}
	if err := s.dao.Transaction(func(tx *gorm.DB) error {
		return s.dao.ActiveUser(ctx, tx, in.Code)
	}); err != nil {
		if s.dao.IsEntityNotFoundError(err) {
			return s.NewError(errInvalidActiveCode)
		}
		return s.InternalServerError(err.Error())
	}
	return nil
}

func (s *userService) VerifyToken(ctx context.Context, in *pb.VerifyTokenRequest, out *pb.VerifyTokenReply) error {
	if len(in.Token) == 0 {
		return s.NewError(errInvalidAccessToken)
	}
	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(in.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.TokenSecret), nil
	})
	if err != nil {
		return s.InternalServerError(err.Error())
	}
	if !token.Valid {
		return s.NewError(errInvalidAccessToken)
	}
	uid, err := strconv.ParseUint(claims.Id, 10, 0)
	if err != nil {
		return s.NewError(errInvalidAccessToken)
	}
	out.UserId = uint32(uid)
	return nil
}

func (s *userService) FellowUser(ctx context.Context, in *pb.FellowUserRequest, out *pb.FellowUserReply) error {
	if err := s.dao.FellowUser(ctx, in.UserId, in.FellowUserId); err != nil {
		return s.InternalServerError(err.Error())
	}
	return nil
}

func (s *userService) GetFellowCount(ctx context.Context, in *pb.GetFellowCountRequest, out *pb.GetFellowCountReply) error {
	fellowCount, err := s.dao.GetFellowUserCount(ctx, in.UserId)
	if err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	}
	fellowerCount, err := s.dao.GetFellowerCount(ctx, in.UserId)
	if err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	}
	out.FellowCount = fellowCount
	out.FellowerCount = fellowerCount
	return nil
}

func (s *userService) GetFellowers(ctx context.Context, in *pb.GetFellowersRequest, out *pb.GetFellowersReply) error {
	total, err := s.dao.GetFellowerCount(ctx, in.UserId)
	if err != nil {
		log.Error(err)
		return s.InternalServerError(err.Error())
	}
	out.Total = total

	return nil
}

func (s *userService) GetFellowedUsers(ctx context.Context, in *pb.GetFellowedUsersRequest, out *pb.GetFellowedUsersReply) error {
	return nil
}

func (s *userService) sendActivationEmail(user *models.User) error {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      "X Community",
			Link:      s.cfg.SiteURL,
			Copyright: fmt.Sprintf("Copyright © %d X Community. All rights reserved.", time.Now().UTC().Year()),
		},
	}
	email := hermes.Email{
		Body: hermes.Body{
			Name:   user.Username,
			Intros: []string{"Welcome to X Community!", "We're very excited to have you on board."},
			Actions: []hermes.Action{{
				Instructions: "To get started with X Community, please click here:",
				Button: hermes.Button{
					Color: "#22BC66",
					Text:  "Confirm your account",
					Link:  s.cfg.SiteURL + "/#/account/active?code=" + user.ActiveCode,
				},
			}},
		},
	}
	subject := "Welcome to X Community"
	content, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}
	_, err = s.mailService.SendMail(context.Background(), &pb.SendMailRequest{Receiver: user.Username, Address: user.Email, Subject: subject, Content: content})
	return err
}
