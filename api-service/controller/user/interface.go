package user

import (
	"github.com/gin-gonic/gin"
	pb "github.com/x-community/api-service/proto"
)

type Service interface {
	Register(*gin.Context)
	VerifyAccount(*gin.Context)
	Login(*gin.Context)
}

func NewService(userService pb.UserService) Service {
	return &service{userService: userService}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email,max=256"`
	Username string `json:"username" binding:"required,max=32"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	EmailOrUsername string `json:"emailOrUsername" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type verifyAccountRequest struct {
	Code string `json:"code" binding:"required,code"`
}
