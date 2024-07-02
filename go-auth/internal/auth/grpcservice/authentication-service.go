package grpcservice

import (
	"context"
	"fmt"
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"gitee.com/go-nianxi/go-auth/internal/pkg/ecode"
	auth "gitee.com/go-nianxi/go-auth/pkg/auth-client/api/authentication"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/client-config"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"go.uber.org/zap"
	"time"
)

type AuthService struct {
	logger    *zap.SugaredLogger
	userDao   *dao.UserDao
	cache     cache.ICache
	jwt       *client_config.Jwt
	jwtOption *client_config.JwtOption

	auth.UnimplementedAuthenticationServer
}

func NewAuthService(userDao *dao.UserDao, appConf *config.AppConf) *AuthService {
	var authService = new(AuthService)
	authService.logger = appConf.Logger.Sugar()
	authService.cache = appConf.Cache
	authService.userDao = userDao
	authService.jwt = appConf.Jwt
	authService.jwtOption = appConf.O.Auth.Jwt

	return authService
}

func (s *AuthService) GetCurrentUserInfo(ctx context.Context, req *auth.UserInfoRequest) (*auth.UserInfoResponse, error) {
	token := req.GetToken()

	resp := &auth.UserInfoResponse{}
	resp.Code = ecode.RECODE_OK

	//验证token字符串
	claim, err := s.jwt.ParserToken(token)
	if err != nil {
		resp.Code = auth_common.ERROR_TOKEN_VALIDATE
		resp.Msg = err.Error()
		return resp, fmt.Errorf("%v\n", err)
	}

	//过期判断
	if time.Now().Unix() > claim.ExpiresAt {
		resp.Code = auth_common.ERROR_TOKEN_EXPIRED
		resp.Msg = auth_common.GetCodeMsg(int(resp.Code))
		return resp, nil
	}

	val, found := s.cache.Get(auth_common.GetAuthCacheKeyString(claim.Id))
	userAuth := val.(auth_common.AuthUser)
	if found {
		userInfo := auth_common.ToUserInfo(&userAuth)
		resp.Data = userInfo
		return resp, nil
	}
	
	resp.Code = ecode.RECODE_USER_NOT_EXISTS
	resp.Msg = ecode.GetCodeMsg(int(resp.Code))
	return resp, nil
}

func (s *AuthService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.BaseResponse, error) {
	resp := &auth.BaseResponse{}
	resp.Code = ecode.RECODE_OK

	token := req.Token

	//验证token字符串
	claim, err := s.jwt.ParserToken(token)
	if err != nil {
		return resp, nil
	}

	//过期判断
	if time.Now().Unix() > claim.ExpiresAt {
		return resp, nil
	}
	if s.cache != nil {
		s.cache.Delete(auth_common.GetAuthCacheKeyString(claim.Id))
	}
	return resp, nil
}
