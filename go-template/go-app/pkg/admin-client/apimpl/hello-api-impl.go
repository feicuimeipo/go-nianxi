package apimpl

import (
	"context"
	"gitee.com/go-nianxi/go-template/pkg/admin-client/api/hello"
	"gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"go.uber.org/zap"
)

type HelloApiImpl struct {
	logger  *zap.Logger
	options *grpc.ClientOptions
}

func NewHelloApiImpl(logger *zap.Logger, options *grpc.ClientOptions) *HelloApiImpl {
	return &HelloApiImpl{
		logger:  logger,
		options: options,
	}
}

func (s *HelloApiImpl) Ping() {
	conn, err := GetConnect(s.options)
	if err != nil {
		s.logger.Fatal("could not greet: %v", zap.Error(err))
		return
	}
	defer conn.Close()

	helloClient := hello.NewHelloClient(conn)

	resp, err := helloClient.Ping(context.Background(), &hello.HelloRequest{Msg: "hello"})
	if err != nil {
		s.logger.Fatal("could not greet: %v", zap.Error(err))
		return
	} else {
		s.logger.Info("---congratulations! auth-server connect successful! A message from auth-server is: %s ", zap.String("resp.msg", resp.Msg))
	}
}
