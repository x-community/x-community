package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/x-community/mail-service/config"
	"github.com/x-community/mail-service/handler"
	tracing "github.com/x-punch/micro-opentracing/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	opts := []micro.Option{micro.Name(cfg.Name), micro.Address(cfg.Address), micro.Version(cfg.Version)}
	if cfg.Tracing.Enable {
		tracer, closer, err := jaegercfg.Configuration{
			ServiceName: cfg.Name,
			Reporter:    &jaegercfg.ReporterConfig{LogSpans: true, CollectorEndpoint: cfg.Tracing.Collector},
		}.NewTracer()
		if err != nil {
			log.Fatal(err)
		}
		defer closer.Close()
		opts = append(opts, micro.WrapHandler(tracing.NewHandlerWrapper(tracer)))
	}
	service := micro.NewService(opts...)
	service.Init()
	if err := handler.Register(service.Server(), handler.Options{ServiceName: cfg.Name, Config: cfg.Handler}); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
