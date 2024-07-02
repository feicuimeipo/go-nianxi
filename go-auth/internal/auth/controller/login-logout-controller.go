package controller

import (
	"errors"
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"gitee.com/go-nianxi/go-auth/internal/auth/domain/vo"
	"gitee.com/go-nianxi/go-auth/internal/auth/model"
	"gitee.com/go-nianxi/go-auth/internal/pkg/ecode"
	auth_client "gitee.com/go-nianxi/go-auth/pkg/auth-client"
	auth_common "gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-common/pkg/captcha"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-common/pkg/utils"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type LoginLogoutController struct {
	valid      *valid.Validator
	userDao    *dao.UserDao
	captcha    *captcha.AJCaptcha
	cache      cache.ICache
	authOption *config.AuthOptions
	authClient *auth_client.ClientContext
	o          *config.AppOptions
}

func NewLoginLogoutController(dao *dao.Dao, appConf *config.AppConf) *LoginLogoutController {
	login := new(LoginLogoutController)
	login.valid = appConf.Valid
	login.userDao = dao.UserDao
	login.captcha = appConf.Captcha
	login.cache = appConf.Cache
	login.authOption = appConf.O.Auth
	login.authClient = appConf.AuthClient
	login.o = appConf.O
	return login
}

// @Tags		认证
// @summary	登录-帐号与密码
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Param		req	body		vo.LoginRequest	true	"登录信息"
// @Router		/auth/login [post]
func (u *LoginLogoutController) Login(c *gin.Context) {

	var req vo.AccountPasswordLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入用户名或密码")
		return
	}

	_, err := u.captcha.HandleFunc.GetCheckCaptchaResult(req.Captcha)
	if err != nil {
		resp.Fail(c, ecode.RECODE_VERIFY_CODE_ERR, "验证码错误")
		return
	}

	if req.LoginName == "" || req.Password == "" {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入用户名或密码")
		return
	}

	if len(req.Password) > 20 {
		decodePassword, err1 := http.RSADecrypt([]byte(req.Password), u.o.Http.Ssl.KeyFileBytes)
		if err1 != nil {
			resp.Fail(c, ecode.RECODE_ECRYPT_ERR, err1.Error())
			return
		}
		req.Password = string(decodePassword)
	}

	var user *model.User
	user, err = u.userDao.GetUserForLogin(req.LoginName, dao.AccountLoginType_USER)
	if err != nil {
		resp.Error(c, ecode.RECODE_USER_NOT_EXISTS, err)
	}

	err = utils.ComparePasswd(user.Password, req.Password)
	if err != nil {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "密码错误！")
		return
	}

	LoginSuccessReturn(c, user, u.authClient)
	return
}

// @Tags		认证
// @summary	登录-手机验证码
// @Accept		json
// @Produce	json
// @Response	200			{object}	resp.ResponseMsg
// @Param		register	body		vo.SMSVerifyCodeLoginRequest	true	"登录信息"
// @Router		/auth/smsLogin [post]
func (u *LoginLogoutController) SMSLogin(c *gin.Context) {
	var req vo.MobileSmsCodeLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入正确的手机号与手机短信验证码")
		return
	}

	var mobile = req.Mobile
	var smsCode = req.SmsCode
	key := config.GetCacheKey(config.RedisPrefixMobileVerifyCode, []string{mobile})
	saveSMSCode, found := u.cache.Get(key)
	if mobile == "" || !u.valid.CheckMobile(mobile) || !found || saveSMSCode != smsCode {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入正确的手机号与手机短信验证码")
		return
	}

	user, err := u.userDao.GetUserForLogin(mobile, dao.AccountLoginType_Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//user = &model.User{Mobile: mobile}
			//user, err = u.userDao.CreateUser(user)
			//if err != nil {
			resp.Fail(c, ecode.RECODE_USER_NOT_EXISTS, "用户不存在，请先注册！")
			//}
		} else {
			resp.Error(c, ecode.RECODE_AUTH_ERR, err)
		}
	}

	LoginSuccessReturn(c, user, u.authClient)
	return
}

func LoginSuccessReturn(c *gin.Context, user *model.User, authClient *auth_client.ClientContext) {
	// 判断用户的状态
	if user.Status != 1 {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "用户被禁用！")
		return
	}

	returnData, err := authClient.LoginSuccess(*model.ToAuthUser(user), c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_JWT_ERR, err.Error())
		return
	}

	resp.Success(c, returnData)
	return
}

func (l *LoginLogoutController) WxQRLogin(c *gin.Context) {
	panic("implement me")
}

func (l *LoginLogoutController) Logout(c *gin.Context) {
	token := c.GetHeader(auth_common.TokenHeaderName)
	if token == "" {
		token, _ = c.Cookie(auth_common.TokenCookiesName)
	}
	c.SetCookie(auth_common.TokenCookiesName, token, time.Now().Second(), "/", l.authOption.Domain, false, true)

	l.authClient.Logout(token)

	resp.Success(c, nil)
}

func (u *LoginLogoutController) ME(c *gin.Context) {
	// 获取当前用户
	ctxUser, err := u.authClient.GetCurrentUser(c)
	if err == nil {
		resp.Success(c, ctxUser)
	} else {
		//resp.Success(c, auth_common.AuthUser{ID: 1, Mobile: "18601106333"})
		resp.Success(c, nil)
	}
}
