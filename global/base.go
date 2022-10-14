package global

import (
	"uni-minds.com/liuxy/medical-sys/logger"
	"uni-minds.com/liuxy/medical-sys/tools"
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
