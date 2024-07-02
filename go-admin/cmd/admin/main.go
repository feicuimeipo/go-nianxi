package main

import (
	"flag"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/etcd"
	"gitee.com/go-nianxi/go-common/pkg/jaeger"
	"gitee.com/go-nianxi/go-admin/internal/admin"
	"gitee.com/go-nianxi/go-admin/internal/admin/config"
	"gitee.com/go-nianxi/go-admin/internal/admin/controller"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/grpcservice"
	"gitee.com/go-nianxi/go-admin/internal/admin/router"
)

var configFile = flag.String("f", "admin.yml", "配置文件路径")

func main() {
	flag.Parsed()

	fmt.Printf("argument number is: %v\n", configFile)

	//初始化配
	viper := base.InitViper(*configFile)
	baseApp := base.NewAppBase(viper)
	app := config.NewApp(baseApp)

	db := config.NewDB(viper, baseApp.Logger)
	casbin := config.NewCasbinEnforcer(db, app.O.Casbin, baseApp.Logger)
	config.InitData(db, base.Conf.Logger)

	config.NewCache(viper, baseApp.Logger)
	valid := config.NewValidator(viper, baseApp.Logger)

	dao := dao.New(db, casbin)
	controller := controller.New(dao, valid, casbin, baseApp.Logger)

	etcd := etcd.NewEtcd(baseApp.Viper, baseApp.Logger)
	jaeger := jaeger.NewJaeger(baseApp.Viper, baseApp.Logger)

	registerService := grpcservice.RegisterService()
	gs, err := grpcservice.NewGrpcServer(baseApp, registerService, jaeger, etcd, baseApp.Logger)
	if err != nil {
		panic(err)
	}

	routers := router.NewRouter(dao, controller, casbin, baseApp.Logger)
	hs := router.NewHttpServer(baseApp, routers, jaeger, etcd)

	application, err := admin.NewApplication(baseApp.O.Name, hs, gs, baseApp.Logger)
	if err != nil {
		panic(err)
	}

	if err := application.Start(); err != nil {
		panic(err)
	}

	application.AwaitSignal()

}
