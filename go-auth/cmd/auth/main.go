package main

import (
	"flag"
	"fmt"
	"gitee.com/go-nianxi/go-auth/internal/auth"
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/controller"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"gitee.com/go-nianxi/go-auth/internal/auth/grpcservice"
	"gitee.com/go-nianxi/go-auth/internal/auth/router"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/etcd"
	"gitee.com/go-nianxi/go-common/pkg/jaeger"
)

var configFile = flag.String("f", "auth.yml", "配置文件路径")

func main() {
	flag.Parsed()

	fmt.Printf("argument number is: %v\n", configFile)

	//初始化配
	viper := base.InitViper(*configFile)
	baseApp := base.NewAppBase(viper)

	valid := config.NewValidator(viper, baseApp.Logger)
	cache := config.NewCache(viper, baseApp.Logger)
	sms := config.NewSMS(viper, baseApp.Logger)
	email := config.NewEmail(viper, baseApp.Logger)
	authClient := config.NewAuthClient(viper)
	captcha := config.NewCaptcha(viper, baseApp.Logger, cache)
	appConf := config.NewConf(baseApp, cache, valid, sms, email, authClient, captcha)

	//cache

	db := config.NewDB(viper, baseApp.Logger)
	dao := dao.New(db)
	controller := controller.New(dao, appConf)

	etcd := etcd.NewEtcd(baseApp.Viper, baseApp.Logger)
	jaeger := jaeger.NewJaeger(baseApp.Viper, baseApp.Logger)

	registerService := grpcservice.RegisterService(dao, appConf)
	gs, err := grpcservice.NewGrpcServer(baseApp, registerService, jaeger, etcd)
	if err != nil {
		panic(err)
	}

	routers := router.NewRouter(dao, controller, appConf)
	hs := router.NewHttpServer(baseApp, routers, jaeger, etcd)

	application, err := auth.NewApplication(baseApp.O.Name, hs, gs, baseApp.Logger)
	if err != nil {
		panic(err)
	}

	if err = application.Start(); err != nil {
		panic(err)
	}

	application.AwaitSignal()

}
