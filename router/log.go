package router

import (
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/medical-sys/tools"
)

const tag = "ROUT"

func log(level string, message ...interface{}) {
	msg := tools.ExpandInterface(message)
	logger.Write(tag, level, msg)
}
