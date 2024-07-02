package controller

import (
	"errors"
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"gitee.com/go-nianxi/go-auth/internal/auth/domain/dto"
	"gitee.com/go-nianxi/go-auth/internal/auth/domain/vo"
	"gitee.com/go-nianxi/go-auth/internal/auth/model"
	"gitee.com/go-nianxi/go-auth/internal/pkg/ecode"
	auth_client "gitee.com/go-nianxi/go-auth/pkg/auth-client"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-common/pkg/utils"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterController struct {
	valid      *valid.Validator
	userDao    *dao.UserDao
	cache      cache.ICache
	authClient *auth_client.ClientContext
	o          *config.AppOptions
	//appConf    *config.AppConf
}

func NewRegisterController(dao *dao.Dao, appConf *config.AppConf) *RegisterController {
	s := new(RegisterController)
	s.valid = appConf.Valid
	s.userDao = dao.UserDao
	s.cache = appConf.Cache
	s.authClient = appConf.AuthClient
	s.o = appConf.O
	//s.appConf = appConf
	return s
}

// @Tags		认证
// @summary	注册
// @Accept		json
// @Produce	json
// @Response	200			{object}	resp.ResponseMsg
// @Param		register	body		vo.RegisterRequest	true	"注册信息"
// @Router		/auth/register/1 [post]
func (u *RegisterController) RegisterStep1(c *gin.Context) {
	var req vo.RegisterStep1Request

	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, u.valid.Translate(err))
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

	users, err := u.userDao.GetUserListLimit5(&model.User{Mobile: req.Mobile})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	} else if err == nil && len(users) > 0 {
		//resp.Success(c, dto.UserIdDTO{ID: users[0].ID, Mobile: users[0].Mobile})
		//return
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "手机号已存在")
		return
	}

	user, err := u.userDao.CreateUser(&model.User{Mobile: mobile})
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	resp.Success(c, dto.UserIdDTO{ID: user.ID, Mobile: req.Mobile})
	return
}

// @Tags		认证
// @summary	注册
// @Accept		json
// @Produce	json
// @Response	200			{object}	resp.ResponseMsg
// @Param		register	body		vo.RegisterRequest	true	"注册信息"
// @Router		/auth/register/2 [post]
func (u *RegisterController) RegisterStep2(c *gin.Context) {
	var req vo.RegisterStep2Request

	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, u.valid.Translate(err))
		return
	}

	if req.Password != req.RePassword {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, "密码与确认密码不相等！")
		return
	}

	users, err := u.userDao.GetUserListLimit5(&model.User{Username: req.UserName})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	} else if err == nil && (len(users) > 1 || (len(users) == 1 && users[0].ID != req.UserId)) {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户名已存在！")
		return
	}

	if len(req.Password) > 20 {
		decodePassword, err1 := http.RSADecrypt([]byte(req.Password), u.o.Http.Ssl.KeyFileBytes)
		if err1 != nil {
			resp.Error(c, ecode.RECODE_ECRYPT_ERR, err1)
			return
		}
		req.Password = string(decodePassword)
	}

	user, err := u.userDao.GetUserById(req.UserId)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	if user.Mobile != req.Mobile {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "手机号不正确！")
		return
	}

	user.Username = req.UserName
	user.Password = utils.GenPasswd(req.Password)
	err = u.userDao.UpdateUser(user)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	LoginSuccessReturn(c, user, u.authClient)
	return
}
