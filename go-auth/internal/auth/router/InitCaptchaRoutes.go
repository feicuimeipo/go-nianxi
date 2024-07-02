package router

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/controller"
	"github.com/gin-gonic/gin"
)

// 安装相关的路帐
func InitCaptchaRoutes(router *gin.RouterGroup, controller *controller.CaptchaController) gin.IRoutes {

	router.POST("/captcha/get", controller.GetCaptcha)
	router.POST("/captcha/check", controller.CheckCaptcha)

	return router
}
