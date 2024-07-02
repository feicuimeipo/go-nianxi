package router

import (
	"gitee.com/go-nianxi/go-admin/internal/admin/controller"
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup, middleware *middleware.Middleware, menuController *controller.MenuController) gin.IRoutes {
	router := r.Group("/menu")
	// 开启jwt认证中间件
	router.Use(middleware.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree/:appId", menuController.GetMenuTree)
		router.GET("/list/:appId", menuController.GetMenus)
		router.POST("/create", menuController.CreateMenu)
		router.PATCH("/update/:menuId", menuController.UpdateMenuById)
		router.DELETE("/delete/batch", menuController.BatchDeleteMenuByIds)
		router.GET("/access/list/:userId", menuController.GetUserMenusByUserId)
		router.GET("/access/tree/:userId/:appId", menuController.GetUserMenuTreeByUserId)
	}

	return r
}
