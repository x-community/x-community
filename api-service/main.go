package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/x-community/api-service/config"
	"github.com/x-community/api-service/controller"
	_ "github.com/x-community/api-service/docs"
	pb "github.com/x-community/api-service/proto"
	mopentracing "github.com/x-punch/micro-opentracing/v2"
)

// @title X Community API Service
// @version 1.0.0
// @description API service for X Community.
// @termsOfService http://swagger.io/terms/
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	api := gin.New()
	opts := []web.Option{web.Name(cfg.Name),
		web.Address(cfg.Address),
		web.Version(cfg.Version),
		web.Handler(api)}
	microOpts := []micro.Option{}
	if cfg.Tracing.Enable {
		jc := jaegercfg.Configuration{
			ServiceName: cfg.Name,
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:          true,
				CollectorEndpoint: cfg.Tracing.Collector,
			},
		}
		tracer, closer, err := jc.NewTracer()
		if err != nil {
			log.Fatal(err)
		}
		defer closer.Close()
		microOpts = append(microOpts, micro.WrapHandler(mopentracing.NewHandlerWrapper(tracer)))
	}
	opts = append(opts, web.MicroService(micro.NewService(microOpts...)))
	srv := web.NewService(opts...)
	controller.Register(controller.Options{Config: cfg.API, Router: api, UserService: pb.NewUserService(cfg.Services.UserService, srv.Options().Service.Client())})
	if err := srv.Init(); err != nil {
		log.Fatal(err)
	}
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
