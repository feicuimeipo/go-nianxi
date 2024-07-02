package router

import (
	"embed"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	_ "gitee.com/go-nianxi/go-admin/api/openapi/admin"
	"gitee.com/go-nianxi/go-admin/internal/admin/controller"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"gitee.com/go-nianxi/go-admin/internal/admin/web"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Router struct {
	HttpRouter     http.InitRouters
	TemplateRouter http.InitTemplateRouter
	staticFs       *embed.FS
	Middleware     *middleware.Middleware
}

func NewRouter(dao *dao.Dao, controller *controller.Controller, casbinEnforcer *casbin.Enforcer, logger *zap.Logger) *Router {
	middleware := middleware.New(dao.ApiDao, dao.UserDao, casbinEnforcer, logger)

	router := http.InitRouters(func(r *gin.RouterGroup, c *gin.Engine) {
		r.Use(middleware.OperationLogMiddleware()) // 启用操作日志中间件

		c.GET("/api-doc/*any", func(c *gin.Context) {
			swagger.DisablingWrapHandler(swaggerFiles.Handler, "SwaggerDisabled")(c)
		})

		// 路由分组
		apiGroup := r.Group(base.UrlPathPrefix)
		// 注册路由
		InitBaseRoutes(apiGroup, middleware)                                            // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
		InitUserRoutes(apiGroup, middleware, controller.UserController)                 // 注册用户路由, jwt认证中间件,casbin鉴权中间件
		InitRoleRoutes(apiGroup, middleware, controller.RoleController)                 // 注册角色路由, jwt认证中间件,casbin鉴权中间件
		InitMenuRoutes(apiGroup, middleware, controller.MenuController)                 // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
		InitApiRoutes(apiGroup, middleware, controller.ApiController)                   // 注册接口路由, jwt认证中间件,casbin鉴权中间件
		InitOperationLogRoutes(apiGroup, middleware, controller.OperationLogController) // 注册操作日志路由, jwt认证中间件,casbin鉴权中间件
		InitAppRoutes(apiGroup, middleware, controller.ApplicationController)
		
	})

	myRouter := Router{
		Middleware:     middleware,
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
