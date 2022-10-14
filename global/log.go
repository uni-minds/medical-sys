package global

import (
	"fmt"
	"github.com/fatih/color"
)

func Warn(str string) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println(str)
	color.Unset()
}

func Error(str string) {
	color.Set(color.FgRed, color.Bold)
	fmt.Println(str)
	color.Unset()
}

func Info(str string) {
	color.Set(color.FgGreen, color.Bold)
	fmt.Println(str)
	color.Unset()
}
