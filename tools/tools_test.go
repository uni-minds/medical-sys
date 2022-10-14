package tools

import (
	"testing"
)

func Test_getMD5(t *testing.T) {
	t.Log(GetFileMD5("a.ogv"))
}

func TestGetStringMD5(t *testing.T) {
	t.Log(GetStringMD5("B101-2020010215040506-2020010216040506.mp4"))
}

func TestGenSaltString(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Log(GenSaltString(20))
	}
}
