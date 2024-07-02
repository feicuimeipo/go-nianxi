package service

import (
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-template/internal/xxx/dao"
	"github.com/google/wire"
)

type Service struct {
	HelloService *HelloService
}

func New(cache cache.ICache, dao *dao.Dao, valid *valid.Validator) *Service {
	return &Service{
		HelloService: NewHelloService(dao, cache, valid),
	}
}

var ProviderSet = wire.NewSet(New)
