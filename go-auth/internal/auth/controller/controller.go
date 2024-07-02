package controller

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"github.com/google/wire"
)

type Controller struct {
	VerifyController         *VerifyController
	UserController           *UserController
	RegisterController       *RegisterController
	LoginLogoutController    *LoginLogoutController
	ForgetPasswordController *ForgetPasswordController
	CaptchaController        *CaptchaController
}

func New(dao *dao.Dao, appConf *config.AppConf) *Controller {
	var controller = new(Controller)
	controller.VerifyController = NewVerifyController(dao, appConf)
	controller.UserController = NewUserController(dao, appConf)
	controller.RegisterController = NewRegisterController(dao, appConf)
	controller.LoginLogoutController = NewLoginLogoutController(dao, appConf)
	controller.CaptchaController = NewCaptchaController(appConf.Captcha)
	controller.ForgetPasswordController = NewForgetPasswordController(dao, appConf)
	return controller
}

var ProviderSet = wire.NewSet(New)
