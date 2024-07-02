package grpc

import (
	"fmt"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	zap_middleware "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	recovery_middleware "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	tags_middleware "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	prometheus_middleware "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type IDiscovery interface {
	register(addr string, serverName string) error
	deRegister(addr string, serverName string) error
	client() interface{}
}

type ServerOptions struct {
	Host          string `mapstructure:"host"              json:"host"`          //端口 >999开启
	Port          int    `mapstructure:"port"              json:"port"`          //端口 >999开启
	CertFile      string `mapstructure:"cert-file"         json:"certFile"`      //证书
	KeyFile       string `mapstructure:"key-file"          json:"keyFile"`       //秘钥
	TrustedCAFile string `mapstructure:"trusted-ca-file"   json:"trustedCAFile"` //ca证书
	Discovery     bool   `mapstructure:"discovery"         json:"discovery"`     //服务发布
}

func NewServerOptions(v *viper.Viper, logger *zap.Logger) *ServerOptions {
	var (
		err error
		o   = new(ServerOptions)
	)
	if err = v.UnmarshalKey("grpc", o); err != nil {
		logger.Sugar().Panicf("初始化配置失败:%s \n", err)
		return nil
	}

	return o
}

type InitServers func(s *grpc.Server)

type Server struct {
	o         *ServerOptions
	host      string
	port      int
	app       string
	logger    *zap.Logger
	server    *grpc.Server
	discovery IDiscovery
	init      InitServers
}

func NewServer(option *ServerOptions, logger *zap.Logger, init InitServers, tracer opentracing.Tracer, etcd *clientv3.Client) (*Server, error) {
	creds, err := tlsServer(option.CertFile, option.KeyFile, option.TrustedCAFile)
	if err != nil {
		logger.Fatal("failed to credentials: %v", zap.Error(err))
		return nil, err
	}
	prometheus_middleware.EnableHandlingTimeHistogram()
	var grpcServer *grpc.Server
	logger = logger.With(zap.String("type", "grpc"))
	grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.StreamInterceptor(
			middleware.ChainStreamServer(
				tags_middleware.StreamServerInterceptor(),
				prometheus_middleware.StreamServerInterceptor,
				zap_middleware.StreamServerInterceptor(logger),
				recovery_middleware.StreamServerInterceptor(),
				otgrpc.OpenTracingStreamServerInterceptor(tracer),
			)),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(
			tags_middleware.UnaryServerInterceptor(),
			prometheus_middleware.UnaryServerInterceptor,
			zap_middleware.UnaryServerInterceptor(logger),
			recovery_middleware.UnaryServerInterceptor(),
			otgrpc.OpenTracingServerInterceptor(tracer),
		)),
	)
	init(grpcServer)

	myServer := &Server{
		o:      option,
		port:   option.Port,
		server: grpcServer,
		logger: logger.With(zap.String("type", "grpc.Server")),
		init:   init,
	}
	if option.Discovery {
		myServer.discovery = NewEtcd(5, etcd)
	}
	return myServer, nil

}

func (s *Server) Start() error {
	s.port = s.o.Port
	s.host = s.o.Host
	if s.port == 0 {
		return errors.New("端口不可以为空！") //s.port = netutil.GetAvailablePort()
	}
	if s.discovery != nil {
		if s.host == "" {
			return errors.New("Host地址不可以为空！")
		}
	}

	addr := fmt.Sprintf("%s:%d", "", s.port)
	s.logger.Info("grpc server starting ...", zap.String("addr", addr))
	go func() {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Fatal("failed to listen: %v", zap.Error(err))
		}

		if err = s.server.Serve(lis); err != nil {
			s.logger.Fatal("failed to serve: %v", zap.Error(err))
		}

	}()

	if s.o.Discovery {
		addr = fmt.Sprintf("%s:%d", s.host, s.port)
		if err := s.discovery.register(addr, s.app); err != nil {
			return errors.Wrap(err, "register grpc server error")
		}
	}

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("grpc server stopping ...")
	if s.discovery != nil {
		addr := fmt.Sprintf("%s:%d", s.host, s.port)
		if err := s.discovery.deRegister(addr, s.app); err != nil {
			return errors.Wrap(err, "deregister grpc server error")
		}
	}
	s.server.GracefulStop()
	return nil
}

func (s *Server) Application(name string) {
	s.app = name
}
