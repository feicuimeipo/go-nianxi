package router

import (
	"gitee.com/go-nianxi/go-admin/internal/admin/controller"
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup, middleware *middleware.Middleware, apiController *controller.ApiController) gin.IRoutes {
	router := r.Group("/api")
	// 开启jwt认证中间件
	router.Use(middleware.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", apiController.GetApis)
		router.GET("/tree/:appId", apiController.GetApiTree)
		router.POST("/create", apiController.CreateApi)
		router.PATCH("/update/:apiId", apiController.UpdateApiById)
		router.DELETE("/delete/batch", apiController.BatchDeleteApiByIds)
	}

	return r
}
