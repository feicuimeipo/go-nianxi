package vo

// 1. 第一步
type ForgetPasswordStep1Request struct {
	Mobile  string `form:"mobile"       json:"mobile"    binding:"required"`
	SmsCode string `form:"smsCode"      json:"smsCode"   binding:"required"`
}

// 2. 重置密码
type ForgetPasswordStep2Request struct {
	NewPassword   string `form:"newPassWord"   json:"newPassword"   binding:"required"`
	ReNewPassword string `form:"reNewPassWord" json:"reNewPassword" binding:"required"`
	Mobile        string `form:"mobile"       json:"mobile"         binding:"required"`
	Username      string `form:"username"     json:"username"`
}
