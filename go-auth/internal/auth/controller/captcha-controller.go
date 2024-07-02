package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/captcha"
	"gitee.com/go-nianxi/go-common/pkg/captcha/core"
	handfunc "gitee.com/go-nianxi/go-common/pkg/captcha/router"
	"github.com/gin-gonic/gin"
)

type CaptchaController struct {
	factory    *core.CaptchaFactory
	HandleFunc *handfunc.HandleFunc
}

func NewCaptchaController(captcha *captcha.AJCaptcha) *CaptchaController {
	s := new(CaptchaController)
	s.factory = captcha.Factory
	s.HandleFunc = captcha.HandleFunc
	return s
}

// @Tags	认证
// @summary	获得验证码
// @Accept		json
// @Produce	json
// @Response	200				{object}	resp.ResponseMsg
// @Param		clientParams	body		handfunc.ClientParams	true	"登录信息"
// @Router		/captcha/get [POST]
func (handler *CaptchaController) GetCaptcha(c *gin.Context) {

	writer := c.Writer
	request := c.Request

	handler.HandleFunc.GetCaptcha(writer, request)
	return

}

// @Tags	认证
// @summary	验证验证码正识
// @Accept		json
// @Produce	json
// @Response	200				{object}	resp.ResponseMsg
// @Param		clientParams	body		handfunc.ClientParams	true	"登录信息"
// @Router		/captcha/check [POST]
func (handler *CaptchaController) CheckCaptcha(c *gin.Context) {
	writer := c.Writer
	request := c.Request

	handler.HandleFunc.CheckCaptcha(writer, request)
	return
}
