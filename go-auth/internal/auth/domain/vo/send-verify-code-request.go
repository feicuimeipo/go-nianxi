package vo

import (
	handfunc "gitee.com/go-nianxi/go-common/pkg/captcha/router"
)

type SendSMSVerifyCodeRequest struct {
	Mobile  string                 `form:"mobile"      json:"mobile"      binding:"required"  example:"18611111111"`
	Use     string                 `form:"use"      json:"use"       example:""` //用途
	Captcha *handfunc.ClientParams `form:"captcha"      json:"captcha"      binding:"required"  example:"18611111111"`
}

type SendEmailVerifyCodeRequest struct {
	Email   string                 `form:"email"      json:"email"      binding:"required"  example:"18611111111"`
	Use     string                 `form:"use"        json:"use"       example:""` //用途
	Captcha *handfunc.ClientParams `form:"captcha"      json:"captcha"      binding:"required"  example:"18611111111"`
}
