package router

import (
	"embed"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	_ "gitee.com/go-nianxi/go-template/api/openapi/admin"
	"gitee.com/go-nianxi/go-template/internal/admin/config"
	"gitee.com/go-nianxi/go-template/internal/admin/controller"
	"gitee.com/go-nianxi/go-template/internal/admin/repository"
	"gitee.com/go-nianxi/go-template/internal/admin/router/middleware"
	"gitee.com/go-nianxi/go-template/internal/admin/web"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var authMiddleware *jwt.GinJWTMiddleware
var casbinMiddleware *middleware.CasbinMiddleWare

type Router struct {
	HttpRouter     http.InitRouters
	TemplateRouter http.InitTemplateRouter
	staticFs       *embed.FS
}

func InitRouter(Repository *repository.Repository, controller *controller.Controller, casbinEnforcer *casbin.Enforcer, logger *zap.Logger) *Router {
	var err error
	authMiddleware, err = middleware.NewAuthMiddleware(config.Conf.O.Http, Repository.SystemRepository)
	if err != nil {
		config.Conf.Logger.Panic("初始化JWT中间件失败", zap.Error(err))
		panic(fmt.Sprintf("初始化JWT中间件失败：%v", err))
	}

	casbinMiddleware = middleware.NewCasbinMiddleWare(authMiddleware, casbinEnforcer, Repository.SystemRepository)

	logMiddle := middleware.NewOperationLogChanMiddleware()

	router := http.InitRouters(func(r *gin.RouterGroup, c *gin.Engine) {
		r.Use(logMiddle.OperationLogMiddleware()) // 启用操作日志中间件

		c.GET("/api-doc/*any", func(c *gin.Context) {
			swagger.DisablingWrapHandler(swaggerFiles.Handler, "SwaggerDisabled")(c)
		})

		// 路由分组
		apiGroup := r.Group(base.UrlPathPrefix)
		InitHelloRoutes(apiGroup, controller.HelloController)
	})

	myRouter := Router{

		HttpRouter:     router,
		TemplateRouter: nil,
		staticFs:       &web.StaticFs,
	}
	return &myRouter

}

func NewHttpServer(appBase *base.BaseConf, router *Router, tracer opentracing.Tracer, etcd *clientv3.Client) *http.Server {
	//构建路由
	httpOptions := http.NewOption(appBase.Viper, appBase.Logger)
	engine := http.NewRouter(httpOptions, router.HttpRouter, tracer, appBase.Logger)
	staticEngines := http.NewStaticRouter(httpOptions, router.TemplateRouter, router.staticFs, appBase.Logger, tracer, engine)
	hs := http.NewServer(httpOptions, engine, staticEngines, etcd, appBase.Logger)
	return hs
}

var ProviderSet = wire.NewSet(InitRouter, NewHttpServer)
