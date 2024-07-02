package grpcservice

import (
	"context"
	"fmt"
	"gitee.com/go-nianxi/go-template/pkg/xxx-client/api/hello"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type HelloService struct {
	logger *zap.SugaredLogger
	hello.UnimplementedHelloServer
}

func NewHelloService(log *zap.SugaredLogger) *HelloService {
	server := new(HelloService)
	server.logger = log
	return server
}

func (h *HelloService) Ping(ctx context.Context, req *hello.HelloRequest) (*hello.BaseResponse, error) {
	// 解析meta_data中的信息并验证
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	for k, v := range md {
		fmt.Println(k, v)
	}

	h.logger.Info("收到来自客户端的ping请求 %s", req.Msg)
	return &hello.BaseResponse{Msg: "pong", Code: 0}, nil
}
