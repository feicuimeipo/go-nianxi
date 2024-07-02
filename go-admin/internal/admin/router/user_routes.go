package router

import (
	"gitee.com/go-nianxi/go-admin/internal/admin/controller"
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"github.com/gin-gonic/gin"
)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup, middleware *middleware.Middleware, userController *controller.UserController) gin.IRoutes {
	router := r.Group("/user")
	// 开启jwt认证中间件
	router.Use(middleware.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/info", userController.GetUserInfo)
		router.GET("/list", userController.GetUsers)
		router.PUT("/changePwd", userController.ChangePwd)
		router.POST("/create", userController.CreateUser)
		router.PATCH("/update/:userId", userController.UpdateUserById)
		router.DELETE("/delete/batch", userController.BatchDeleteUserByIds)
	}
	return r
}
