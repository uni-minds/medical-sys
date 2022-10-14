/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: base.go
 */

package global

import (
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/medical-sys/tools"
)

const tag = "GLOB"

var debugMode bool

func log(level string, message ...interface{}) {
	msg := tools.ExpandInterface(message)
	logger.Write(tag, level, msg)
}

func DebugSetFlag(f bool) {
	if f {
		log("w", "Debug Mode")
		debugMode = true
	} else {
		debugMode = false
	}
}

func DebugGetFlag() bool {
	return debugMode
}
