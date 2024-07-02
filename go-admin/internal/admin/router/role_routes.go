package router

import (
	"gitee.com/go-nianxi/go-admin/internal/admin/controller"
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoleRoutes(r *gin.RouterGroup, middleware *middleware.Middleware, roleController *controller.RoleController) gin.IRoutes {
	router := r.Group("/role")
	// 开启jwt认证中间件
	router.Use(middleware.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", roleController.GetRoles)
		router.GET("/all", roleController.GetAllRoles)
		router.GET("/menus/get/:roleId", roleController.GetRoleMenusById)
		router.GET("/apis/get/:roleId", roleController.GetRoleApisById)
		router.POST("/create", roleController.CreateRole)
		router.PATCH("/update/:roleId", roleController.UpdateRoleById)
		router.PATCH("/menus/update/:roleId", roleController.UpdateRoleMenusById)
		router.PATCH("/apis/update/:roleId", roleController.UpdateRoleApisById)
		router.DELETE("/delete/batch", roleController.BatchDeleteRoleByIds)
	}
	return r
}
