package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/x-community/x-community/user-service/config"
	"github.com/x-community/x-community/user-service/database"
	"github.com/x-community/x-community/user-service/handler"
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
	if err := handler.RegisterHandlers(service.Server(), handler.Options{ServiceName: cfg.Name, DB: db}); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
