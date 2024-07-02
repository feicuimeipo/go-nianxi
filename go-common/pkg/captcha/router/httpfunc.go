package router

import (
	"encoding/json"
	"errors"
	"gitee.com/go-nianxi/go-common/pkg/captcha/core"
	"io/ioutil"
	"net/http"
)

type HandleFunc struct {
	factory *core.CaptchaFactory
}

func NewHandleFunc(factory *core.CaptchaFactory) *HandleFunc {
	return &HandleFunc{factory: factory}
}

type ClientParams struct {
	Token       string `json:"token"           example:"1234"`
	PointJson   string `json:"pointJson"       example:"1234"`
	CaptchaType string `json:"captchaType"     example:"1234"`
}

/**
 * 行为校验配置模块（具体参数可从业务系统配置文件自定义）
 */
func (handler *HandleFunc) Cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                      // 可将将 * 替换为指定的域名
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,x-requested-with") //你想放行的header也可以在后面自行添加
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")                                                    //我自己只使用 get post 所以只放行它
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// 放行所有OPTIONS方法
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// 处理请求
		f(w, r)
	}
}

func (handler *HandleFunc) GetCaptcha(writer http.ResponseWriter, request *http.Request) {
	params, err := handler.GetParams(request)
	if err != nil {
		res, _ := json.Marshal(handler.Error(err))
		writer.Write(res)
		return
	}
	if params.CaptchaType == "" {
		res, _ := json.Marshal(handler.Error(errors.New("参数传递不完整")))
		writer.Write(res)
		return
	}

	ser := handler.factory.GetCaptchaService(params.CaptchaType)

	data, err := ser.Get()

	res, _ := json.Marshal(handler.Success(data))
	request.Body.Close()
	writer.Write(res)
}

func (handler *HandleFunc) CheckCaptcha(writer http.ResponseWriter, request *http.Request) {
	params, err := handler.GetParams(request)

	if params.Token == "" || params.PointJson == "" || params.CaptchaType == "" {
		res, _ := json.Marshal(handler.Error(errors.New("参数传递不完整")))
		writer.Write(res)
		return
	}

	if err != nil {
		res, _ := json.Marshal(handler.Error(err))
		writer.Write(res)
		return
	}

	ser := handler.factory.GetCaptchaService(params.CaptchaType)

	err = ser.Check(params.Token, params.PointJson)
	if err != nil {
		res, _ := json.Marshal(handler.Error(err))
		writer.Write(res)
		return
	}

	res, _ := json.Marshal(handler.Success(nil))
	writer.Write(res)
}

/**
 * 得到验证结果
 */
func (handler *HandleFunc) GetCheckCaptchaResult(params *ClientParams) ([]byte, error) {

	//params, err := handler.GetParams(request)

	if params == nil || params.Token == "" || params.PointJson == "" || params.CaptchaType == "" {
		res, _ := json.Marshal(handler.Error(errors.New("参数传递不完整")))
		return res, errors.New("参数传递不完整")
	}

	ser := handler.factory.GetCaptchaService(params.CaptchaType)

	err := ser.Check(params.Token, params.PointJson)
	if err != nil {
		res, _ := json.Marshal(handler.Error(err))
		return res, err
	}
	res, _ := json.Marshal(handler.Success(nil))
	return res, nil
}

func (handler *HandleFunc) GetParams(request *http.Request) (*ClientParams, error) {
	params := &ClientParams{}
	all, _ := ioutil.ReadAll(request.Body)

	if len(all) <= 0 {
		query := request.URL.Query()
		params.CaptchaType = query.Get("captchaType")
		params.PointJson = query.Get("pointJson")
		params.Token = query.Get("token")

		if params.CaptchaType == "" && params.PointJson == "" && params.Token == "" {
			return nil, errors.New("未获取到客户端数据")
		}
	}

	err := json.Unmarshal(all, params)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (handler *HandleFunc) Error(err error) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["code"] = "0001"
	ret["data"] = nil
	ret["msg"] = err.Error()
	return ret
}

func (handler *HandleFunc) Success(data interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["code"] = "0"
	ret["msg"] = ""
	ret["data"] = data
	return ret
}
