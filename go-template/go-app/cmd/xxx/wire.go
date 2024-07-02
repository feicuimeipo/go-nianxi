//go:build wireinject
// +build wireinject

package main

import (
	"gitee.com/go-nianxi/go-common/pkg/app"
	"gitee.com/go-nianxi/go-template/internal/xxx"
	"gitee.com/go-nianxi/go-template/internal/xxx/config"
	"gitee.com/go-nianxi/go-template/internal/xxx/dao"
	"gitee.com/go-nianxi/go-template/internal/xxx/grpcservice"
	"gitee.com/go-nianxi/go-template/internal/xxx/router"
	"gitee.com/go-nianxi/go-template/internal/xxx/service"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	config.ProviderSet,
	dao.ProviderSet,
	service.ProviderSet,
	grpcservice.ProviderSet,
	router.ProviderSet,
	xxx.ProviderSet,
)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
