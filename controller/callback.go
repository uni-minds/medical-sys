/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: callback.go
 * Description:
 */

package controller

import "net/http"

func FailReturn(code int, msg interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	res["data"] = ""
	res["code"] = code
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
