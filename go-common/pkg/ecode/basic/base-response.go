package basic

/**
 * 通用链码返回值
 */
type BasicResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(data interface{}) *BasicResponse {
	return &BasicResponse{RECODE_OK, "OK", data}
}

func ERR(code int, msg string) *BasicResponse {
	return &BasicResponse{code, msg, nil}
}
