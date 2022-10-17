package logger

import (
	"fmt"
	"gitee.com/uni-minds/utils/tools"
)

type Logger struct {
	tag string
}

func NewLogger(tag string) *Logger {
	return &Logger{tag: tag}
}

func (logger *Logger) Printf(format string, a ...any) {
	Write(logger.tag, "i", fmt.Sprintf(format, a...))
}

func (logger *Logger) Println(contents ...any) {
	Write(logger.tag, "i", tools.InterfaceExpand(contents))
}

func (logger *Logger) Error(content string) {
	Write(logger.tag, "e", content)
}
func (logger *Logger) Errorf(format string, a ...any) {
	Write(logger.tag, "e", fmt.Sprintf(format, a...))
}

func (logger *Logger) Debug(content string) {
	Write(logger.tag, "d", content)
}

func (logger *Logger) Debugf(format string, a ...any) {
	Write(logger.tag, "d", fmt.Sprintf(format, a...))
}

func (logger *Logger) Trace(content string) {
	Write(logger.tag, "t", content)
}
func (logger *Logger) Warn(content string) {
	Write(logger.tag, "w", content)
}

func (logger *Logger) Warnf(format string, a ...any) {
	Write(logger.tag, "w", fmt.Sprintf(format, a...))
}

func (logger *Logger) Log(level string, message ...interface{}) {
	msg := tools.InterfaceExpand(message)
	Write(logger.tag, level, msg)
}
