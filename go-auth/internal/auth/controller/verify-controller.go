package controller

import (
	"fmt"
	"gitee.com/go-nianxi/go-auth/internal/auth/config"
	"gitee.com/go-nianxi/go-auth/internal/auth/dao"
	"gitee.com/go-nianxi/go-auth/internal/auth/domain/vo"
	"gitee.com/go-nianxi/go-auth/internal/pkg/ecode"
	"gitee.com/go-nianxi/go-auth/internal/pkg/email"
	"gitee.com/go-nianxi/go-auth/internal/pkg/sms"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-common/pkg/captcha"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"time"
)

var VerifyCache = freecache.NewCache(100 * 1024 * 1024)

type VerifyController struct {
	cache     cache.ICache
	sms       *sms.SMSClient
	email     *email.MailClient
	valid     *valid.Validator
	ajCaptcha *captcha.AJCaptcha
	userDao   *dao.UserDao
	logger    *zap.Logger
	o         *config.AppOptions
}

func NewVerifyController(dao *dao.Dao, appConf *config.AppConf) *VerifyController {
	s := new(VerifyController)
	s.userDao = dao.UserDao
	s.cache = appConf.Cache
	s.valid = appConf.Valid
	s.sms = appConf.Sms
	s.ajCaptcha = appConf.Captcha
	s.logger = appConf.Logger
	s.o = appConf.O
	s.email = appConf.Email
	return s
}

// @Tags		认证
// @summary		发送手机验证码
// @description	use:1-注册帐号时确认 2:短信验证码登录 3:安全绑定 ...
// @Accept			json
// @Produce		json
// @Response		200	{object}	resp.ResponseMsg
// @Param			req	body		vo.SendSMSVerifyCodeRequest	true	"发送手机验证码的请求数据"
// @Router			/auth/sendSmsVerifyCode [post]
func (controller *VerifyController) SendSmsVerifyCode(c *gin.Context) {
	var req vo.SendSMSVerifyCodeRequest
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, controller.valid.Translate(err))
		return
	}

	//判断验证码,验证通过
	_, err := controller.ajCaptcha.HandleFunc.GetCheckCaptchaResult(req.Captcha)
	if err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, err.Error())
		return
	}

	var mobile = req.Mobile

	if mobile == "" || !controller.valid.CheckMobile(mobile) {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, "请输入正确的手机号！")
		return
	}

	key := config.GetCacheKey(config.RedisPrefixMobileVerifyCode, []string{req.Mobile, req.Use})
	exist := controller.cache.Exists(key)
	if exist {
		resp.Fail(c, ecode.RECODE_VERIFY_CODE_ERR, fmt.Sprintf("%d秒后再试！", 60))
		return
	}
	digital := CreateCode()
	controller.cache.Set(key, digital, time.Second*60)

	//v, boo := controller.cache.Get(key)
	//config.Conf.Logger.Info("缓存的数据,", zap.String("value", v.(string)), zap.Bool("bool", boo))

	controller.logger.Sugar().Info("验证码:", key, ":", digital)

	//生产模式
	if !controller.o.Auth.Dev.IsVerifyLocal {
		controller.sms.SendSMS(mobile, fmt.Sprintf("手机验证码：%s", digital))
	}
	resp.Success(c, digital)
}

// @Tags		认证
// @summary		发送手机验证码
// @description	use:1-注册帐号时确认 2:短信验证码登录 3:安全绑定 ...
// @Accept			json
// @Produce		json
// @Response		200	{object}	resp.ResponseMsg
// @Param			req	body		vo.SendEmailVerifyCodeRequest	true	"发送邮箱验证码的请求数据"
// @Router			/auth/sendEmailVerifyCode [post]
func (controller *VerifyController) SendEmailVerifyCode(c *gin.Context) {
	var req vo.SendEmailVerifyCodeRequest
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, controller.valid.Translate(err))
		return
	}

	//判断验证码,验证通过
	_, err := controller.ajCaptcha.HandleFunc.GetCheckCaptchaResult(req.Captcha)
	if err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, err.Error())
		return
	}

	var email = req.Email

	if email == "" || !controller.valid.CheckEmail(email) {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, "请输入正确的邮箱帐号！")
		return
	}

	key := config.GetCacheKey(config.RedisPrefixEmailVerifyCode, []string{req.Email, req.Use})
	exist := controller.cache.Exists(key)
	if exist {
		resp.Fail(c, ecode.RECODE_VERIFY_CODE_ERR, fmt.Sprintf("%d秒后再试！", 60))
		return
	}
	digital := CreateCode()
	controller.cache.Set(key, digital, time.Second*60)

	controller.logger.Sugar().Info("验证码:", digital)

	//生产模式
	if !controller.o.Auth.Dev.IsVerifyLocal {
		controller.email.SendMail(email, "验证码", fmt.Sprintf("验证码：%s", digital))
	}
	resp.Success(c, digital)
}

func CreateCode() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Int31n(10000)) //这里面前面的04v是和后面的1000相对应的
}
