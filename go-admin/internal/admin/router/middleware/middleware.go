package middleware

import (
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-admin/internal/admin/config"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

var middleware *Middleware

type Middleware struct {
	apiDao         *dao.ApiDao
	userDao        *dao.UserDao
	casbin         *casbin.Enforcer
	AuthMiddleware *jwt.GinJWTMiddleware
	logger         *zap.Logger
	httOption      *http.Options
}

func New(apiDao *dao.ApiDao, userDao *dao.UserDao, casbin *casbin.Enforcer, logger *zap.Logger) *Middleware {
	middleware = &Middleware{}
	middleware.apiDao = apiDao
	middleware.userDao = userDao
	middleware.casbin = casbin
	middleware.logger = logger
	middleware.httOption = config.Conf.O.Http
	var err error
	middleware.AuthMiddleware, err = middleware.InitAuth()
	if err != nil {
		config.Conf.Logger.Panic("初始化JWT中间件失败", zap.Error(err))
		panic(fmt.Sprintf("初始化JWT中间件失败：%v", err))
	}
	return middleware
}
