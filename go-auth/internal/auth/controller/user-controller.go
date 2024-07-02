package controller

import (
	"errors"
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
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

type UserController struct {
	cache      cache.ICache
	userDao    *dao.UserDao
	valid      *valid.Validator
	authClient *auth_client.ClientContext
	o          *config.AppOptions
}

func NewUserController(dao *dao.Dao, appConf *config.AppConf) *UserController {
	s := new(UserController)
	s.valid = appConf.Valid
	s.cache = appConf.Cache
	s.userDao = dao.UserDao
	s.authClient = appConf.AuthClient
	s.o = appConf.O
	return s
}

// @Tags	用户
// @summary	更新用户名
// @Accept		json
// @Produce	json
// @Param		register	body		vo.ChangeUsernameRequest	true	"修改密码"
// @Response	200			{object}	resp.ResponseMsg
// @Router		/user/update/username [PATCH]
func (u *UserController) UpdateUserName(c *gin.Context) {
	var req vo.ChangeUsernameRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	// 获取当前用户
	ctxUser, err := u.authClient.GetCurrentUser(c)
	if err != nil || (ctxUser != nil && ctxUser.ID != req.UserId) {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	users, err := u.userDao.GetUserListLimit5(&model.User{Username: req.UserName})
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}
	if len(users) > 1 || (len(users) == 1 && users[0].ID != ctxUser.ID) {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "用户名已存在！")
		return
	}

	user, err1 := u.userDao.GetUserById(req.UserId)
	if err1 != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}
	user.Username = req.UserName
	err = u.userDao.UpdateUser(user)

	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	authUser := model.ToAuthUser(user)

	resp.Success(c, authUser)
	return
}

// @Tags	用户
// @summary	更新昵称
// @Accept		json
// @Produce	json
// @Param		register	body		vo.ChangeNickNameRequest	true	"修改昵称"
// @Response	200			{object}	resp.ResponseMsg
// @Router		/user/update/nickname [PATCH]
func (u *UserController) UpdateNickname(c *gin.Context) {
	var req vo.ChangeNickNameRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	if req.NickName == "" {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "请输入正确的昵称")
		return
	}

	// 获取当前用户
	ctxUser, err := u.authClient.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	if ctxUser.ID != req.UserId {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "不是当前用户！")
		return
	}

	users, err := u.userDao.GetUserListLimit5(&model.User{Nickname: req.NickName})
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}
	if len(users) > 1 || (len(users) == 1 && users[0].ID != ctxUser.ID) {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "已存在！")
		return
	}

	user, err := u.userDao.GetUserById(req.UserId)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	user.Nickname = req.NickName
	err = u.userDao.UpdateUser(user)

	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	authUser := model.ToAuthUser(user)

	resp.Success(c, authUser)
	return
}

// @Tags	用户
// @summary	更新昵称
// @Accept		json
// @Produce	json
// @Param		register	body		vo.ChangeNickNameRequest	true	"修改昵称"
// @Response	200			{object}	resp.ResponseMsg
// @Router		/user/update/mobile [PATCH]
func (u *UserController) UpdateMobile(c *gin.Context) {
	var req vo.ChangeMobileRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	key := config.GetCacheKey(config.RedisPrefixMobileVerifyCode, []string{req.NewMobile, "update"})
	oriValue, found := u.cache.Get(key)
	if req.NewMobile == "" || !u.valid.CheckMobile(req.NewMobile) || !found || req.SmsCode != oriValue {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入正确的手机号与手机短信验证码")
		return
	}

	// 获取当前用户
	ctxUser, err := u.authClient.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	if ctxUser.ID != req.UserId {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "不是当前用户！")
		return
	}

	users, err := u.userDao.GetUserListLimit5(&model.User{Mobile: req.NewMobile})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}
	if len(users) > 1 || (len(users) == 1 && users[0].ID != ctxUser.ID) {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "手机号已存在！")
		return
	}

	user, err1 := u.userDao.GetUserById(req.UserId)
	if err1 != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	user.Mobile = req.NewMobile
	err = u.userDao.UpdateUser(user)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	authUser := model.ToAuthUser(user)

	resp.Success(c, authUser)
	return
}

func (u *UserController) UpdateEmail(c *gin.Context) {
	var req vo.ChangeEmailRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	key := config.GetCacheKey(config.RedisPrefixEmailVerifyCode, []string{req.Email, "update"})
	oriValue, found := u.cache.Get(key)
	if req.Email == "" || !u.valid.CheckEmail(req.Email) || !found || req.VerifyCode != oriValue {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "请输入正确的邮箱号码与邮箱验证码")
		return
	}

	// 获取当前用户
	ctxUser, err := u.authClient.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	if ctxUser.ID != req.UserId {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "不是当前用户！")
		return
	}

	users, err := u.userDao.GetUserListLimit5(&model.User{Email: req.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}
	if len(users) > 1 || (len(users) == 1 && users[0].ID != ctxUser.ID) {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "邮箱已存在！")
		return
	}

	user, err1 := u.userDao.GetUserById(req.UserId)
	if err1 != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	user.Email = req.Email
	err = u.userDao.UpdateUser(user)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	authUser := model.ToAuthUser(user)

	resp.Success(c, authUser)
	return
}

// @Tags		用户
// @summary	修改密码
// @Accept		json
// @Produce	json
// @Response	200			{object}	resp.ResponseMsg
// @Param		register	body		vo.ChangePwdRequest	true	"修改密码"
// @Router		/user/changePwd [put]
func (u *UserController) ChangePwd(c *gin.Context) {
	var req vo.ChangePwdRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, err.Error())
		return
	}

	if req.ReNewPassword == "" || req.NewPassword != req.ReNewPassword {
		resp.Fail(c, ecode.RECODE_USER_PASSWORD_ERROR, "新密码与确认密码不对！")
		return
	}

	userInfo, err := u.authClient.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, err.Error())
		return
	}

	if len(req.OldPassword) > 20 {
		decodeOldPassword, err1 := http.RSADecrypt([]byte(req.OldPassword), u.o.Http.Ssl.KeyFileBytes)
		if err1 != nil {
			resp.Fail(c, ecode.RECODE_ECRYPT_ERR, err1.Error())
			return
		}
		decodeNewPassword, err1 := http.RSADecrypt([]byte(req.NewPassword), u.o.Http.Ssl.KeyFileBytes)
		if err1 != nil {
			resp.Error(c, ecode.RECODE_ECRYPT_ERR, err1)
			return
		}
		req.OldPassword = string(decodeOldPassword)
		req.NewPassword = string(decodeNewPassword)
	}

	user, err := u.userDao.GetUserById(userInfo.ID)
	if err != nil {
		resp.Error(c, ecode.RECODE_INTERNAL_ERR, err)
		return
	}
	// 判断前端请求的密码是否等于真实密码
	err = utils.ComparePasswd(user.Password, utils.GenPasswd(req.OldPassword))
	if err != nil {
		resp.Fail(c, ecode.RECODE_AUTH_ERR, "原密码不正确！")
		return
	}

	// 更新密码
	user.Password = utils.GenPasswd(req.NewPassword)
	err = u.userDao.UpdateUser(user)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新密码失败: "+err.Error())
		return
	}
	resp.Success(c, nil)
}
