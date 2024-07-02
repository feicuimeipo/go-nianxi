package apimpl

import (
	"gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"go.uber.org/zap"
	googleGrpc "google.golang.org/grpc"
)

type ApiImpl struct {
	HelloClient *HelloApiImpl
}

func NewApiImpl(options *grpc.ClientOptions, logger *zap.Logger) *ApiImpl {
	return &ApiImpl{
		HelloClient: NewHelloApiImpl(logger, options),
	}
}

func GetConnect(options *grpc.ClientOptions) (*googleGrpc.ClientConn, error) {
	client, err := grpc.NewClient(options, nil, nil)
	if err != nil {
		return nil, err
	}
	conn, err := client.Dial()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
