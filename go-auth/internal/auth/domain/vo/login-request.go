package vo

import handfunc "gitee.com/go-nianxi/go-common/pkg/captcha/router"

// 1.用户密码登录
type AccountPasswordLoginRequest struct {
	LoginName string                 `form:"loginName" json:"loginName" binding:"required"   example:"faker"`
	Password  string                 `form:"password" json:"password"   binding:"required"   example:"123456"`
	Captcha   *handfunc.ClientParams `form:"captcha"      json:"captcha"      binding:"required"  example:"18611111111"`
}

// 1.手机验证码登录
type MobileSmsCodeLoginRequest struct {
	Mobile  string `form:"mobile"     json:"mobile"      binding:"required"          example:"18612345678"`
	SmsCode string `form:"smsCode"     json:"smsCode"    binding:"required"          example:"1234"`
}

// 手机扫_二维码登录
type WxQRCodeLoginRequest struct {
	QRCode   string `form:"qrCode" json:"qrCode"`
	WxOpenId string `form:"wxOpenId"         json:"wxOpenId"`
}
