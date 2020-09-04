package controller

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/x-community/api-service/controller/account"
	"github.com/x-community/api-service/docs"
	pb "github.com/x-community/api-service/proto"
)

// Options represents api service options
type Options struct {
	Config      Config
	Router      *gin.Engine
	UserService pb.UserService
}

func Register(opts Options) error {
	if opts.Config.Swagger.Enable {
		docs.SwaggerInfo.Host = opts.Config.Swagger.Host
		docs.SwaggerInfo.BasePath = opts.Config.Swagger.Base
		opts.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router := opts.Router
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	accountService := account.NewService(opts.UserService)
	{
		router.POST("/account/register", accountService.Register)
		router.POST("/account/verify", accountService.VerifyAccount)
		router.POST("/account/login", accountService.Login)
	}
	return nil
}
