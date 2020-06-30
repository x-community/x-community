package user

import (
	"github.com/gin-gonic/gin"
	"github.com/x-community/api-service/controller/api"
	pb "github.com/x-community/api-service/proto"
)

var _ Service = &service{}

type service struct {
	api.APIService
	userService pb.UserService
}

// @Description register
// @Tags Account
// @ID account_register
// @Param	input 	body		user.registerRequest		true		"request"
// @Success 200 	{object}	api.SuccessResponse	"success"
// @Failure 400 	{object}	api.ErrorResponse
// @Failure 500		{object}	api.ErrorResponse
// @Router /account/register [post]
func (s *service) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBind(&req); nil != err {
		s.BadRequest(ctx, api.InvalidRequest(err.Error()))
		return
	}
	_, err := s.userService.Register(s.TracerContext(ctx), &pb.RegisterRequest{Email: req.Email, Username: req.Username, Password: req.Password})
	if err != nil {
		s.Error(ctx, err)
		return
	}
	s.Success(ctx)
}

// @Description verify account
// @Tags Account
// @ID account_verify
// @Param	input 	body		user.verifyAccountRequest	true		"request"
// @Success 200 	{object}	api.SuccessResponse			"success"
// @Failure 400 	{object}	api.ErrorResponse
// @Failure 500		{object}	api.ErrorResponse
// @Router /account/verify [post]
func (s *service) VerifyAccount(ctx *gin.Context) {
	var req verifyAccountRequest
	if err := ctx.ShouldBind(&req); nil != err {
		s.BadRequest(ctx, api.InvalidRequest(err.Error()))
		return
	}
	_, err := s.userService.VerifyAccount(s.TracerContext(ctx), &pb.VerifyAccountRequest{Code: req.Code})
	if err != nil {
		s.Error(ctx, err)
		return
	}
	s.Success(ctx)
}

// @Security ApiKeyAuth
// @Description login
// @Tags Account
// @ID account_login
// @Param	input 	body		user.loginRequest	true	"request"
// @Success 200 	{object}	user.loginResponse			"success"
// @Failure 401 	{object}	api.ErrorResponse
// @Failure 500		{object}	api.ErrorResponse
// @Router /account/logout [post]
func (s *service) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBind(&req); nil != err {
		s.BadRequest(ctx, api.InvalidRequest(err.Error()))
		return
	}
	resp, err := s.userService.Authenticate(s.TracerContext(ctx), &pb.AuthenticateRequest{EmailOrUsername: req.EmailOrUsername, Password: req.Password})
	if err != nil {
		s.Error(ctx, err)
		return
	}
	s.Success(ctx, loginResponse{Token: resp.Token})
}
