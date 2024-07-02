package apimpl

import (
	"context"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/api/authentication"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/api/hello"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	nianxiGrpc "gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"go.uber.org/zap"
)

type AuthApiImpl struct {
	grpc   *nianxiGrpc.ClientOptions
	logger *zap.SugaredLogger
}

func newAuthApiClientImp(grpc *nianxiGrpc.ClientOptions, logger *zap.Logger) *AuthApiImpl {
	return &AuthApiImpl{
		grpc:   grpc,
		logger: logger.Sugar(),
	}
}

func (s *AuthApiImpl) Logout(ctx context.Context, token string) (*authentication.BaseResponse, error) {
	conn, err := getConnect(s.grpc)
	if err != nil {
		s.logger.Fatal(err)
		return nil, err
	}
	defer conn.Close()
	in := authentication.LogoutRequest{}
	in.Token = token
	authClient := authentication.NewAuthenticationClient(conn)
	resp, err := authClient.Logout(ctx, &in)
	return resp, err
}

func (s *AuthApiImpl) TryingConnect() {
	conn, err := getConnect(s.grpc)
	if err != nil {
		s.logger.Fatalf("could not greet: %v", err)
		return
	}
	defer conn.Close()

	helloClient := hello.NewHelloClient(conn)

	resp, err := helloClient.Ping(context.Background(), &hello.HelloRequest{Msg: "hello"})
	if err != nil {
		s.logger.Fatalf("could not greet: %v", err)
		return
	}
	s.logger.Info("---congratulations! auth-server connect successful! A message from auth-server is: %s ", resp.Msg)
}

func (s *AuthApiImpl) GetCurrentUserInfo(ctx context.Context, token string) (*auth_common.AuthUser, error) {
	conn, err := getConnect(s.grpc)
	if err != nil {
		s.logger.Fatal(err)
		return nil, err
	}
	defer conn.Close()

	in := authentication.UserInfoRequest{}
	in.Token = token
	authClient := authentication.NewAuthenticationClient(conn)
	resp, err := authClient.GetCurrentUserInfo(ctx, &in)
	if err != nil {
		s.logger.Panic("认证服务器出错！", err)
	}

	dto := auth_common.ToAuthUser(resp.GetData())

	return dto, nil
}
