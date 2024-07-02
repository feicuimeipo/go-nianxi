package valid

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"regexp"
)

// CheckMobile 检验手机号
func (v *Validator) CheckMobile(phone string) bool {
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}

// CheckIdCard 检验身份证
func (v *Validator) CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}

// 识别电子邮箱
func (v *Validator) CheckEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)
	if result {
		return true
	} else {
		return false
	}
}

func (v *Validator) ParseParams(request *http.Request, params interface{}) error {
	all, _ := ioutil.ReadAll(request.Body)

	if len(all) <= 0 {
		return errors.New("未获取到客户端数据")
	}

	err := json.Unmarshal(all, params)
	if err != nil {
		return err
	}

	err = v.Validate.Struct(params)
	if err != nil {
		return errors.New(v.Translate(err))
	}

	return nil
}

func (v *Validator) Translate(err error) (ret string) {
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return err.Error()
	} else {
		for _, e := range validationErrors {
			ret += e.Translate(v.Trans) + ";"
		}
	}
	return ret
}
