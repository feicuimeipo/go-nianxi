package ecode

import "fmt"

const (
	RECODE_OK           = 0
	RECODE_BadRequest   = 400
	RECODE_AUTH_JWT_ERR = 401
	RECODE_AUTH_ERR     = 402

	RECODE_ECRYPT_ERR      = 403
	RECODE_VERIFY_CODE_ERR = 404

	RECODE_INTERNAL_ERR  = 500
	RECODE_REQ_PARAM_ERR = 501
	RECODE_LimitVisit    = 502
	RECODE_UNKNOWN_ERR   = 999
	RECODE_JSON_ERR      = 904

	//区块链相用相关
	RECODE_CC_PUT_ERR   = 911
	RECODE_CC_GET_ERR   = 912
	RECODE_CC_DEL_ERR   = 913
	RECODE_INVOKE_ERR   = 914
	RECODE_CC_CALL_ERR  = 906
	RECODE_CC_GETID_ERR = 917

	//fabric-用户相关
	RECODE_EXIST_MULTI_USERS   = 800
	RECODE_USER_ALREADY_EXISTS = 801
	RECODE_USER_NOT_EXISTS     = 802
	RECODE_USER_PASSWORD_ERROR = 803
	RECODE_GET_USER_ERR        = 804

	//fabric-token
	RECODE_BALANCE_NOT_ENOUGH = 701

	//fabric-task
	RECODE_ZONEID_GET_ERR             = 622
	RECODE_TASK_STATUS_ERR            = 640
	RECODE_TASK_NOT_EXIST             = 642
	RECODE_TASK_CANNOT_TAKE           = 643
	RECODE_TASK_CANNOT_COMMIT         = 644
	RECODE_TASK_CANNOT_CONFIRM        = 645
	RECODE_TASK_CANNOT_BOUNS_TRANSFER = 646
	RECODE_TASK_CANNOT_CANCEL         = 648

	//http-gitee.com/go-nianxi/go-auth-sdk
	RECODE_BIND_ERR = 931
)

var mapRespMsg map[int]string = map[int]string{
	RECODE_OK:           "成功！",
	RECODE_BadRequest:   "请求失败",
	RECODE_LimitVisit:   "限流访问",
	RECODE_INTERNAL_ERR: "系统异常",
	RECODE_UNKNOWN_ERR:  "未知错误",
	RECODE_JSON_ERR:     "JSON操作错误",
	RECODE_AUTH_ERR:     "权限错误",
	RECODE_AUTH_JWT_ERR: "JWT认证失败",

	RECODE_CC_PUT_ERR:   "写区块错误",
	RECODE_CC_GET_ERR:   "读块错误",
	RECODE_CC_DEL_ERR:   "册区块错误",
	RECODE_CC_CALL_ERR:  "链码调用失败",
	RECODE_INVOKE_ERR:   "跨链调用异常",
	RECODE_CC_GETID_ERR: "获得ID有误",

	RECODE_USER_ALREADY_EXISTS: "用户已存在",
	RECODE_USER_NOT_EXISTS:     "用户不存在",
	RECODE_USER_PASSWORD_ERROR: "密码错误",
	RECODE_GET_USER_ERR:        "获取用户信息失败",
	RECODE_BALANCE_NOT_ENOUGH:  "资金不足",

	//fabric-task
	RECODE_ZONEID_GET_ERR:             "获取zoneId错误",
	RECODE_TASK_STATUS_ERR:            "任务状态错误",
	RECODE_TASK_NOT_EXIST:             "任务不存在",
	RECODE_TASK_CANNOT_TAKE:           "任务不能领取",
	RECODE_TASK_CANNOT_COMMIT:         "任务不能提交",
	RECODE_TASK_CANNOT_CONFIRM:        "任务不能确认",
	RECODE_TASK_CANNOT_BOUNS_TRANSFER: "任务转帐异常",
	RECODE_TASK_CANNOT_CANCEL:         "任务不能取消",

	RECODE_BIND_ERR: "数据绑定错误",
}

func GetCodeMsg(code int) string {
	value := mapRespMsg[code]
	if value == "" {
		return fmt.Sprintf("[%d]", code)
	}
	return value
}
