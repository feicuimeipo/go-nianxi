package middleware

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
)

var appConf *config.AppConf

type Middleware struct {
}

func New(dao *dao.Dao) *Middleware {
	var middleware = &Middleware{}
	return middleware
}
