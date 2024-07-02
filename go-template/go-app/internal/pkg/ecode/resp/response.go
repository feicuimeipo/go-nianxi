package resp

import (
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/ecode/basic"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PageData struct {
	List  interface{}
	Total int64
}

type ResponseMsg struct {
	Code int         `json:"code" example:"0"`
	Msg  string      `json:"msg"  example:"code=0，msg=成功，data为返回值，code!=0时，msg返回错误信息"`
	Data interface{} `json:"data" `
}

// message:可空
func Response(c *gin.Context, httpStatus int, code int, data interface{}, message ...string) {
	resp := NewResponseMsg()
	resp.Code = code
	if code == http.StatusOK {
		resp.Code = basic.RECODE_OK
	}

	msg := ""
	for _, v := range message {
		msg = fmt.Sprintf("%s%s\n", msg, v)
	}

	resp.Msg = msg
	resp.Data = data

	c.JSON(httpStatus, gin.H{"data": resp.Data, "msg": resp.Msg, "code": resp.Code})
}

func JSON(c *gin.Context, resp *ResponseMsg) {
	c.JSON(http.StatusOK, resp)
}

func Fail(c *gin.Context, code int, message ...string) {
	Response(c, http.StatusOK, code, nil, message...)
}

func Error(c *gin.Context, code int, err error) {
	Response(c, http.StatusOK, code, nil, err.Error())
}

func OK(c *gin.Context, data interface{}, message ...string) {
	Response(c, http.StatusOK, basic.RECODE_OK, data, message...)
}

func SuccessGinH(c *gin.Context, data gin.H, message ...string) {
	Response(c, http.StatusOK, basic.RECODE_OK, data, message...)
}

func SuccessPageList(c *gin.Context, list interface{}, total int64, message ...string) {
	Response(c, http.StatusOK, basic.RECODE_OK, gin.H{"list": list, "total": total}, message...)
}

func SuccessPage(c *gin.Context, pageData PageData, message ...string) {
	Response(c, http.StatusOK, basic.RECODE_OK, gin.H{"list": pageData.List, "total": pageData.Total}, message...)
}

func Success(c *gin.Context, data interface{}, message ...string) {
	Response(c, http.StatusOK, basic.RECODE_OK, data, message...)
}

func Writer(c *gin.Context, httpStatus int, code int, data interface{}, message ...string) {
	Response(c, httpStatus, code, data, message...)
}

func NewResponseMsg() ResponseMsg {
	return ResponseMsg{basic.RECODE_OK, "", nil}
}
