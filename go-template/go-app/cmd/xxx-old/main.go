package main

import (
	"flag"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/etcd"
	"gitee.com/go-nianxi/go-common/pkg/jaeger"
	"gitee.com/go-nianxi/go-template/internal/xxx"
	"gitee.com/go-nianxi/go-template/internal/xxx/config"
	"gitee.com/go-nianxi/go-template/internal/xxx/dao"
	"gitee.com/go-nianxi/go-template/internal/xxx/grpcservice"
	"gitee.com/go-nianxi/go-template/internal/xxx/router"
	"gitee.com/go-nianxi/go-template/internal/xxx/service"
)

var configFile = flag.String("f", "xxx.yml", "配置文件路径")

func main() {
	flag.Parsed()

	fmt.Printf("argument number is: %v\n", configFile)

	//初始化配
	viper := base.InitViper(*configFile)
	baseApp := base.NewAppBase(viper)
	config.NewApp(baseApp)
	db := config.NewDB(viper, baseApp.Logger)
	cache := config.NewCache(viper, baseApp.Logger)
	valid := config.NewValidator(viper, baseApp.Logger)

	//dao
	dao := dao.New(db)
	service := service.New(cache, dao, valid)

	etcd := etcd.NewEtcd(baseApp.Viper, baseApp.Logger)
	jaeger := jaeger.NewJaeger(baseApp.Viper, baseApp.Logger)

	registerService := grpcservice.RegisterService(service, cache, baseApp.Logger.Sugar())
	gs, err := grpcservice.NewGrpcServer(baseApp, registerService, jaeger, etcd)
	if err != nil {
		panic(err)
	}

	routers := router.NewRouter(service, valid)
	hs := router.NewHttpServer(baseApp, routers, jaeger, etcd)

	application, err := xxx.NewApp(baseApp.O.Name, hs, gs, baseApp.Logger)
	if err != nil {
		panic(err)
	}

	if err := application.Start(); err != nil {
		panic(err)
	}

	application.AwaitSignal()

}
