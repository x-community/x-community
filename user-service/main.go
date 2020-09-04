package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/x-community/user-service/config"
	"github.com/x-community/user-service/database"
	"github.com/x-community/user-service/handler"
	pb "github.com/x-community/user-service/proto"
	tracing "github.com/x-punch/micro-opentracing/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewDatabase(cfg.DB)
	opts := []micro.Option{micro.Name(cfg.Name), micro.Address(cfg.Address), micro.Version(cfg.Version)}
	if cfg.Tracing.Enable {
		cfg := jaegercfg.Configuration{
			ServiceName: cfg.Name,
			Sampler:     &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
			Reporter:    &jaegercfg.ReporterConfig{LogSpans: true, CollectorEndpoint: cfg.Tracing.Collector},
		}
		tracer, closer, err := cfg.NewTracer()
		if err != nil {
			log.Fatal(err)
		}
		defer closer.Close()
		opts = append(opts, micro.WrapHandler(tracing.NewHandlerWrapper(tracer)))
	}
	service := micro.NewService(opts...)
	service.Init()
	mailService := pb.NewMailService(cfg.Services.MailService, service.Client())
	if err := handler.Register(service.Server(), handler.Options{ServiceName: cfg.Name, Config: cfg.Handler, DB: db, MailService: mailService}); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
