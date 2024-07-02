package valid

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"regexp"
)

type Options struct {
	Local string `mapstructure:"language" json:"local"`
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("app", o); err != nil {
		return nil, err
	}

	return o, err
}

type Validator struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

// 初始化Validator数据校验
func New(config *Options, logger *zap.Logger) (*Validator, error) {
	// 全局Validate数据校验实列
	var validate *validator.Validate

	// 全局翻译器
	var trans ut.Translator

	var ok bool
	if validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		local := config.Local
		trans, o = uni.GetTranslator(local)
		if !o {
			logger.Sugar().Error(errors.New("uni.GetTranslator failed"))
			return nil, nil
		}

		// register translate
		// 注册翻译器
		var err error
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(validate, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(validate, trans)
		default:
			err = chTranslations.RegisterDefaultTranslations(validate, trans)
		}

		if err != nil {
			logger.Sugar().Error("初始化验证语言失败", zap.Error(err))
			return nil, nil
		}

	}

	// 添加额外翻译
	_ = validate.RegisterTranslation("required_without", trans, func(ut ut.Translator) error {
		return ut.Add("required_without", "{0} 为必填字段!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_without", fe.Field())
		return t
	})
	_ = validate.RegisterTranslation("mobile_without", trans, func(ut ut.Translator) error {
		return ut.Add("mobile_without", "{0} 手机格式错误!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile_without", fe.Field())
		return t
	})
	_ = validate.RegisterTranslation("memail_without", trans, func(ut ut.Translator) error {
		return ut.Add("email_without", "{0} 邮箱格式错误!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email_without", fe.Field())
		return t
	})

	validate.RegisterValidation("mobile", checkMobile)
	validate.RegisterValidation("email", checkEmail)
	logger.Sugar().Infof("初始化validator.v10数据校验器完成")

	return &Validator{validate, trans}, nil
}

func checkMobile(fl validator.FieldLevel) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}

// 识别电子邮箱
func checkEmail(fl validator.FieldLevel) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, fl.Field().String())
	if result {
		return true
	} else {
		return false
	}
}
