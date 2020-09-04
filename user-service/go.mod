module github.com/x-community/user-service

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/matcornic/hermes/v2 v2.1.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/x-punch/go-config v1.0.5
	github.com/x-punch/go-strings v0.1.1
	github.com/x-punch/micro-opentracing/v2 v2.0.4
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/mysql v1.0.1
	gorm.io/gorm v1.20.0
)
