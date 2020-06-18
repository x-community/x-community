module github.com/x-community/x-community/user-service

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.13
	github.com/micro/go-micro/v2 v2.9.0
	github.com/uber/jaeger-client-go v2.24.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/x-punch/go-config v1.0.5
	github.com/x-punch/micro-opentracing/v2 v2.0.2
	google.golang.org/protobuf v1.24.0
)
