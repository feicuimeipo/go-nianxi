package grpc_simple

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	localgrpc "gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"path/filepath"
)

const schema = "ns"

type GrpcServer struct {
	serverName string
	logger     *zap.Logger
	o          *localgrpc.ServerOptions
	cli        *clientv3.Client
	v          *viper.Viper
	init       localgrpc.InitServers
}

func NewGrpcServer(serverName string, init localgrpc.InitServers, cli *clientv3.Client, v *viper.Viper, logger *zap.Logger) *grpc.Server {
	server := GrpcServer{
		logger: logger,
	}
	o := localgrpc.NewServerOptions(v, logger)
	server.serverName = serverName
	server.cli = cli
	server.v = v
	server.o = o
	server.init = init
	return server.SimpleGrpcStart()
}

func tlsServer(option *localgrpc.ServerOptions) credentials.TransportCredentials {
	crtFile := filepath.Join("key", option.CertFile)
	keyFile := filepath.Join("key", option.KeyFile)
	caFile := filepath.Join("key", option.TrustedCAFile)

	// TLS认证
	//从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	//key, _ := tls.LoadX509KeyPair("certs/server.pem","certs/server.key")
	cert, _ := tls.LoadX509KeyPair(crtFile, keyFile)
	certPool := x509.NewCertPool() //初始化一个CertPool
	ca, _ := os.ReadFile(caFile)
	certPool.AppendCertsFromPEM(ca) //解析传入的证书，解析成功会将其加到池子中
	cred := credentials.NewTLS(&tls.Config{ //构建基于TLS的TransportCredentials选项
		Certificates: []tls.Certificate{cert},        //服务端证书链，可以有多个
		ClientAuth:   tls.RequireAndVerifyClientCert, //要求必须验证客户端证书
		ClientCAs:    certPool,                       //设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})
	return cred
}

func (g *GrpcServer) SimpleGrpcStart() *grpc.Server {

	host := "127.0.0.1"
	port := fmt.Sprintf(":%d", g.o.Port)
	serverAddr := fmt.Sprintf("%s:%d", host, port)

	g.logger.Info("--GRPC服务启动---")

	cred := tlsServer(g.o)
	srv := grpc.NewServer(grpc.Creds(cred))
	//服务注册
	g.init(srv)

	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			g.logger.Sugar().Fatal("failed to listen: %v", err)
		}

		if err = srv.Serve(lis); err != nil {
			g.logger.Error("failed to serve: %v", zap.Error(err))
		}
	}()

	//记动注册中心
	if g.o.Discovery {
		g.logger.Info("--grpc启动etcd注册中心--")
		InitRegisterServer(serverAddr, g.serverName, g.cli)
	}

	g.logger.Sugar().Info("grpc服务启动完成,地址：%s,服务名：%s", serverAddr, g.serverName)
	return srv
}
