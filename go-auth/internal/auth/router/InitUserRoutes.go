package router

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/controller"
	"github.com/gin-gonic/gin"
)

// 用户
func InitUserRoutes(r *gin.RouterGroup, userController *controller.UserController) gin.IRoutes {
	router := r.Group("/user")
	{
		router.PATCH("/update/password", userController.ChangePwd)
		router.PATCH("/update/username", userController.UpdateUserName)
		router.PATCH("/update/nickname", userController.UpdateNickname)
		router.PATCH("/update/mobile", userController.UpdateMobile)
		router.PATCH("/update/email", userController.UpdateEmail)

	}

	return r
}
