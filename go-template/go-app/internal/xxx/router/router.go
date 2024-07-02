package router

import (
	"embed"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	_ "gitee.com/go-nianxi/go-template/api/openapi/xxx"
	"gitee.com/go-nianxi/go-template/internal/xxx/router/middleware"
	"gitee.com/go-nianxi/go-template/internal/xxx/service"
	"gitee.com/go-nianxi/go-template/internal/xxx/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Router struct {
	HttpRouter     http.InitRouters
	TemplateRouter http.InitTemplateRouter
	staticFs       *embed.FS
}

func NewRouter(service *service.Service, valid *valid.Validator) *Router {
	router := http.InitRouters(func(r *gin.RouterGroup, c *gin.Engine) {

		//自定义路由中间件
		c.Use(middleware.HelloRequiredMiddleware(c))

		// 初始化JWT认证中间件
		c.GET("/api-doc/*any", func(c *gin.Context) {
			swagger.DisablingWrapHandler(swaggerFiles.Handler, "SwaggerDisabled")(c)
		})

		// 注册路由
		InitHelloRoutes(r, service, valid) // 权限路由, 不需要jwt认证中间件
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

var ProviderSet = wire.NewSet(NewRouter, NewHttpServer)
