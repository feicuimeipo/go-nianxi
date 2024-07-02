package router

import (
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"github.com/gin-gonic/gin"
)

// 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup, middleware *middleware.Middleware) gin.IRoutes {
	router := r.Group("/auth")
	{
		// 登录登出刷新token无需鉴权
		router.POST("/login", middleware.AuthMiddleware.LoginHandler)
		router.POST("/logout", middleware.AuthMiddleware.LogoutHandler)
		router.POST("/refreshToken", middleware.AuthMiddleware.RefreshHandler)
	}
	return r
}
