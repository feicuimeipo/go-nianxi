package grpcservice

import (
	"gitee.com/go-nianxi/go-template/pkg/admin-client/api/hello"
	"gitee.com/go-nianxi/go-common/pkg/base"
	pkgGrpc "gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"github.com/opentracing/opentracing-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

func RegisterService() pkgGrpc.InitServers {
	grpcServer := pkgGrpc.InitServers(func(s *grpc.Server) {
		hello.RegisterHelloServer(s, NewHelloService())
	})
	return grpcServer
}

func NewGrpcServer(appBase *base.BaseConf, registerServers pkgGrpc.InitServers, jaeger opentracing.Tracer, etcd *clientv3.Client, logger *zap.Logger) (*pkgGrpc.Server, error) {
	grpcOption := pkgGrpc.NewServerOptions(appBase.Viper, logger)
	gs, err := pkgGrpc.NewServer(grpcOption, appBase.Logger, registerServers, jaeger, etcd)
	return gs, err
}

var ProviderSet = wire.NewSet(RegisterService, NewGrpcServer)
