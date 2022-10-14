package module

import (
	"testing"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func TestGenSaltString(t *testing.T) {
	t.Log(tools.GenSaltString(6))
}
