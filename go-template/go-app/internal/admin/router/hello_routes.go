package router

import (
	"gitee.com/go-nianxi/go-template/internal/admin/controller"
	"github.com/gin-gonic/gin"
)

func InitHelloRoutes(r *gin.RouterGroup, hello *controller.HelloController) gin.IRoutes {
	router := r.Group("/hello")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(casbinMiddleware.MiddlewareFunc())
	{
		router.GET("/list", hello.GetHellos)
	}
	return r
}
