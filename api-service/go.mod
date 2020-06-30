module github.com/x-community/api-service

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.9.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/uber/jaeger-client-go v2.24.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/x-punch/gin-logger v1.0.2
	github.com/x-punch/go-config v1.0.5
	github.com/x-punch/micro-opentracing/v2 v2.0.2
	google.golang.org/protobuf v1.25.0
)
