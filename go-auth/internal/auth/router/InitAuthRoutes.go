package router

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/controller"
	"github.com/gin-gonic/gin"
)

// 安装相关的路帐
func InitAuthRoutes(r *gin.RouterGroup, controller *controller.Controller) gin.IRoutes {

	router := r.Group("/auth")
	{
		router.POST("/login", controller.LoginLogoutController.Login)
		router.POST("/smsLogin", controller.LoginLogoutController.SMSLogin)
		router.POST("/register/1", controller.RegisterController.RegisterStep1)
		router.POST("/register/2", controller.RegisterController.RegisterStep2)
		router.POST("/sendSmsVerifyCode", controller.VerifyController.SendSmsVerifyCode)
		router.POST("/sendEmailVerifyCode", controller.VerifyController.SendEmailVerifyCode)
		router.POST("/findPassword/1", controller.ForgetPasswordController.ForgetPasswordStep1)
		router.POST("/findPassword/2", controller.ForgetPasswordController.ForgetPasswordStep2ResetPassword)
		router.GET("/logout", controller.LoginLogoutController.Logout)
		router.GET("/me", controller.LoginLogoutController.ME)
	}
	return r
}
