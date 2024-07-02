package controller

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"gitee.com/go-nianxi/go-auth/internal/auth/domain/dto"
	"gitee.com/go-nianxi/go-auth/internal/auth/domain/vo"
	"gitee.com/go-nianxi/go-auth/internal/auth/model"
	"gitee.com/go-nianxi/go-auth/internal/pkg/ecode"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-common/pkg/captcha"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-common/pkg/utils"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"github.com/gin-gonic/gin"
)

type ForgetPasswordController struct {
	userDao   *dao.UserDao
	valid     *valid.Validator
	ajCaptcha *captcha.AJCaptcha
	cache     cache.ICache
	o         *config.AppOptions
}

func NewForgetPasswordController(dao *dao.Dao, appConf *config.AppConf) *ForgetPasswordController {
	forgetPassword := new(ForgetPasswordController)
	forgetPassword.ajCaptcha = appConf.Captcha
	forgetPassword.userDao = dao.UserDao
	forgetPassword.cache = appConf.Cache
	forgetPassword.valid = appConf.Valid
	forgetPassword.o = appConf.O
	return forgetPassword
}

func (f *ForgetPasswordController) ForgetPasswordStep1(c *gin.Context) {
	var req vo.ForgetPasswordStep1Request
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入手机号与用户名")
		return
	}

	var mobile = req.Mobile
	var smsCode = req.SmsCode
	key := config.GetCacheKey(config.RedisPrefixMobileVerifyCode, []string{mobile})
	saveSMSCode, found := f.cache.Get(key)
	if mobile == "" || !f.valid.CheckMobile(mobile) || !found || saveSMSCode != smsCode {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入正确的手机号与手机短信验证码")
		return
	}

	users, err1 := f.userDao.GetUserListLimit5(&model.User{
		Mobile: req.Mobile,
	})
	if err1 != nil {
		resp.Fail(c, ecode.RECODE_USER_NOT_EXISTS, err1.Error())
		return
	}
	if len(users) == 0 {
		resp.Fail(c, ecode.RECODE_USER_NOT_EXISTS, "用户不存在！")
		return
	}

	multiRecord := false
	if len(users) > 1 {
		multiRecord = true
	}

	resp.Success(c, dto.UserIdDTO{ID: users[0].ID, Mobile: req.Mobile, MultiRecord: multiRecord})
}

func (f *ForgetPasswordController) ForgetPasswordStep2ResetPassword(c *gin.Context) {
	var req vo.ForgetPasswordStep2Request
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "信息不对")
		return
	}
	if req.ReNewPassword == "" && req.NewPassword != req.ReNewPassword {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "重置密码不对！")
		return
	}

	if len(req.NewPassword) > 20 {
		decodePassword, err1 := http.RSADecrypt([]byte(req.ReNewPassword), f.o.Http.Ssl.KeyFileBytes)
		if err1 != nil {
			resp.Fail(c, ecode.RECODE_ECRYPT_ERR, err1.Error())
			return
		}
		req.NewPassword = string(decodePassword)
	}

	search := &model.User{
		Mobile: req.Mobile,
	}
	if req.Username != "" {
		search.Username = req.Username
	}

	users, err1 := f.userDao.GetUserListLimit5(search)

	if err1 != nil {
		resp.Fail(c, ecode.RECODE_USER_NOT_EXISTS, err1.Error())
		return
	}
	if len(users) == 0 || len(users) > 1 {
		resp.Fail(c, ecode.RECODE_USER_NOT_EXISTS, "该用户一个手机绑定多个帐号，请输入正确的用户名？")
		return
	}

	user := users[0]
	user.Password = utils.GenPasswd(req.NewPassword)
	f.userDao.UpdateUser(user)

	resp.Success(c, model.ToAuthUser(user))
}
