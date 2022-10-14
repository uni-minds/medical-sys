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

func TestRemoveDuplicateInt(t *testing.T) {
	data := []int{1, 2, 4, 6, 2}
	o := RemoveDuplicateInt(data)
	t.Log(o)
}

func TestRemoveElementInt(t *testing.T) {
	data := []int{1, 2, 4, 6, 2}
	t.Log(RemoveElementInt(data, 1))
}
