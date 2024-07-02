package router

import (
	"gitee.com/go-nianxi/go-admin/internal/admin/controller"
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"github.com/gin-gonic/gin"
)

func InitAppRoutes(r *gin.RouterGroup, middleware *middleware.Middleware, controller *controller.ApplicationController) gin.IRoutes {
	router := r.Group("/app")
	// 开启jwt认证中间件
	router.Use(middleware.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/type/list", controller.GetApplicationTypes)
		router.GET("/tree", controller.GetApplicationTree)
		router.GET("/list", controller.GetApplications)
		router.POST("/create", controller.CreateApplication)
		router.PATCH("/update/:id", controller.UpdateApplicationById)
		router.DELETE("/delete/batch", controller.BatchDeleteApiByIds)
		router.GET("/get/:userId", controller.GetAppsByUserId)
	}

	return r
}
