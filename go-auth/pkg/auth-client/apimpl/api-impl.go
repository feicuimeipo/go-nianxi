package apimpl

import (
	nianxiGrpc "gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ApiImpl struct {
	AuthApi *AuthApiImpl
}

func NewApiImpl(grpc *nianxiGrpc.ClientOptions, logger *zap.Logger) *ApiImpl {
	return &ApiImpl{
		AuthApi: newAuthApiClientImp(grpc, logger),
	}
}

func getConnect(grpc *nianxiGrpc.ClientOptions) (*grpc.ClientConn, error) {
	client, err := nianxiGrpc.NewClient(grpc, nil, nil)
	if err != nil {
		return nil, err
	}
	conn, err := client.Dial()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
