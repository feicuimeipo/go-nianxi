package service

import (
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-template/internal/xxx/dao"
	"gitee.com/go-nianxi/go-template/internal/xxx/domain/vo"
	"gitee.com/go-nianxi/go-template/internal/xxx/model"
)

type HelloService struct {
	helloDao *dao.HelloDao
	valid    *valid.Validator
	cache    cache.ICache
}

func NewHelloService(dao *dao.Dao, cache cache.ICache, valid *valid.Validator) *HelloService {
	return &HelloService{
		valid:    valid,
		cache:    cache,
		helloDao: dao.HelloDao,
	}
}

func (h *HelloService) GetHelloById(req *vo.HelloRequest) (*model.Hello, error) {
	return h.helloDao.GetHelloById(req.Id)

}
