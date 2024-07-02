package grpcservice

import (
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	pkgGrpc "gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"gitee.com/go-nianxi/go-template/internal/xxx/service"
	"gitee.com/go-nianxi/go-template/pkg/xxx-client/api/hello"
	"github.com/opentracing/opentracing-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

func RegisterService(service *service.Service, cache cache.ICache, log *zap.SugaredLogger) pkgGrpc.InitServers {
	grpcServer := pkgGrpc.InitServers(func(s *grpc.Server) {
		hello.RegisterHelloServer(s, NewHelloService(log))
	})
	return grpcServer
}

func NewGrpcServer(appBase *base.BaseConf, registerServers pkgGrpc.InitServers, jaeger opentracing.Tracer, etcd *clientv3.Client) (*pkgGrpc.Server, error) {
	grpcOption := pkgGrpc.NewServerOptions(appBase.Viper, appBase.Logger)
	gs, err := pkgGrpc.NewServer(grpcOption, appBase.Logger, registerServers, jaeger, etcd)
	return gs, err
}

var ProviderSet = wire.NewSet(RegisterService, NewGrpcServer)
