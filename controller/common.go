package controller

import "net/http"

func FailReturn(msg interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	res["data"] = ""
	res["code"] = http.StatusBadRequest
	res["msg"] = msg

	return res
}

// SuccessReturn api正确返回函数
func SuccessReturn(msg interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	res["data"] = msg
	res["code"] = http.StatusOK
	res["msg"] = "success"

	return res
}

const (
	ETokenInvalid     = "登录凭证无效"
	EActionForbiden   = "禁止操作"
	EParameterInvalid = "参数异常"
)
