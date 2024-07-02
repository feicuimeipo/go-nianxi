package router

import (
	"embed"
	_ "gitee.com/go-nianxi/go-auth/api/openapi/auth"
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/controller"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"gitee.com/go-nianxi/go-auth/internal/auth/router/middleware"
	"gitee.com/go-nianxi/go-auth/internal/auth/web"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Router struct {
	HttpRouter http.InitRouters
	staticFs   *embed.FS
	Middleware *middleware.Middleware
}

func NewRouter(dao *dao.Dao, ctrl *controller.Controller, appConf *config.AppConf) *Router {
	middleware := middleware.New(dao)

	router := http.InitRouters(func(r *gin.RouterGroup, c *gin.Engine) {
		//设置模式
		gin.SetMode(base.Conf.O.Mode)
		gin.ForceConsoleColor()

		r.Use(appConf.AuthClient.NxAuthMiddleware())
		r.Use(middleware.OperationLogMiddleware(appConf.O.Http.ContextPath)) // 启用操作日志中间件
		// 初始化JWT认证中间件

		c.GET("/api-doc/*any", func(c *gin.Context) {
			swagger.DisablingWrapHandler(swaggerFiles.Handler, "SwaggerDisabled")(c)
		})

		// 注册路由
		InitCaptchaRoutes(r, ctrl.CaptchaController) // 注册用户路由, jwt认证中间件
		InitAuthRoutes(r, ctrl)                      // 权限路由, 不需要jwt认证中间件
		InitUserRoutes(r, ctrl.UserController)       // 注册用户路由, jwt认证中间件
		appConf.AuthClient.InitAuthClientRouter(r)
	})

	myRouter := Router{
		Middleware: middleware,
		HttpRouter: router,
		staticFs:   &web.StaticFs,
	}
	return &myRouter
}

func NewHttpServer(appBase *base.BaseConf, router *Router, tracer opentracing.Tracer, etcd *clientv3.Client) *http.Server {
	//构建路由
	httpOptions := http.NewOption(appBase.Viper, appBase.Logger)
	engine := http.NewRouter(httpOptions, router.HttpRouter, tracer, appBase.Logger)
	staticEngines := http.NewStaticRouter(httpOptions, router.staticFs, appBase.Logger, tracer, engine)
	hs := http.NewServer(httpOptions, engine, staticEngines, etcd, appBase.Logger)
	return hs
}

var ProviderSet = wire.NewSet(NewRouter, NewHttpServer)
