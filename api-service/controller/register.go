package controller

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/x-community/api-service/controller/user"
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
	userService := user.NewService(opts.UserService)
	{
		router.POST("/account/register", userService.Register)
		router.POST("/account/verify", userService.VerifyAccount)
		router.POST("/account/login", userService.Login)
	}
	return nil
}
