/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: base.go
 */

package controller

import (
	"uni-minds.com/liuxy/medical-sys/logger"
	"uni-minds.com/liuxy/medical-sys/tools"
)

const tag = "CTRL"
const edaAddress = "localhost:80"       // http
const mboxGateway = "localhost:8442"    // https
const algoServer = "192.168.2.101:5000" // https

func log(level string, message ...interface{}) {
	msg := tools.ExpandInterface(message)
	logger.Write(tag, level, msg)
}

const (
	ETokenInvalid     = "登录凭证无效"
	EActionForbiden   = "禁止操作"
	EParameterInvalid = "参数异常"
)
