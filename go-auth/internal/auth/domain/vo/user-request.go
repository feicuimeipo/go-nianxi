package vo

type RegisterStep1Request struct {
	Mobile  string `form:"mobile"     json:"mobile"        binding:"mobile"     example:"18612345678"`
	SmsCode string `form:"smsCode"    json:"smsCode"       binding:"required"     example:"18612345678"`
}

type RegisterStep2Request struct {
	UserId     uint   `form:"userId"   json:"userId"         binding:"required"   example:"nianxi"`
	UserName   string `form:"username"   json:"username"      binding:"required"`
	Password   string `form:"password"   json:"password"      binding:"required"`
	RePassword string `form:"rePassword" json:"rePassword"    binding:"required"`
	Mobile     string `form:"mobile"     json:"mobile"        binding:"required"     example:"18612345678"`
}

// 修改密码
type ChangePwdRequest struct {
	UserId        uint   `form:"userId"        json:"userId"        binding:"required"    example:"用户编号"`
	OldPassword   string `form:"oldPassword"   json:"oldPassword"   binding:"required"    example:"旧密码"`
	NewPassword   string `form:"newPassword"   json:"newPassword"   binding:"required"    example:"新密码"`
	ReNewPassword string `form:"reNewPassword" json:"reNewPassword" binding:"required"    example:"重复新密码"`
}

type ChangeUsernameRequest struct {
	UserId   uint   `form:"userId"      json:"userId"        binding:"required"    example:"用户编号"`
	UserName string `form:"userName"   json:"userName"       binding:"required"    example:"用户名"`
}

type ChangeMobileRequest struct {
	UserId    uint   `form:"userId"        json:"userId"        binding:"required"    example:"用户编号"`
	NewMobile string `form:"newMobile"     json:"newMobile"        binding:"required"     example:"18612345678"`
	SmsCode   string `form:"smsCode"    json:"smsCode"       binding:"required"     example:"18612345678"`
}

type ChangeEmailRequest struct {
	UserId     uint   `form:"userId"        json:"userId"        binding:"required"    example:"用户编号"`
	Email      string `form:"email"         json:"email"        binding:"required"     example:"18612345678"`
	VerifyCode string `form:"verifyCode"    json:"verifyCode"       binding:"required"     example:"18612345678"`
}

type ChangeNickNameRequest struct {
	UserId   uint   `form:"userId"       json:"userId"        binding:"required"    example:"用户编号"`
	NickName string `form:"nickName"     json:"nickName"        binding:"required"   example:"18612345678"`
}
