/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: base.go
 */

package tools

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/logger"
)

const tag = "TOOL"

func log(level string, message ...interface{}) {
	msg := ExpandInterface(message)
	logger.Write(tag, level, msg)
}

func ExpandInterface(msg []interface{}) string {
	str := fmt.Sprintf("%v", msg[0])
	for _, v := range msg[1:] {
		str = fmt.Sprintf("%s %v", str, v)
	}
	return str
}
