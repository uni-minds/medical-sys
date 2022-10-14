/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: base.go
 */

package global

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/utils/tools"
)

const tag = "GLOB"

var flagDebug, flagVerbose bool

func log(level string, message ...interface{}) {
	logger.Write(tag, level, tools.InterfaceExpand(message))
}

func FlagSetDebug(f bool) {
	if f {
		fmt.Println("Debug On")
	}
	flagDebug = f
}

func FlagGetDebug() bool {
	return flagDebug
}

func FlagSetVerbose(f bool) {
	if f {
		fmt.Println("Verbose On")
	}
	flagVerbose = f
}

func FlagGetVerbose() bool {
	return flagVerbose
}
