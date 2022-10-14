package tools

import (
	"testing"
)

func TestStringCompress(t *testing.T) {
	//data1 := []string{"a","b","c"}
	data1 := []string{}
	str, err := StringCompress(data1)
	t.Log(str, err)

	data2 := ""
	strD, err := StringDecompress(data2)
	t.Log(strD, err)
}

func TestStringsDedup(t *testing.T) {
	data := []string{"a", "b", "1", "d", "c", "b", "c"}
	t.Log(StringsDedup(data))
}

func TestStringsDedupWithSort(t *testing.T) {
	str1 := []string{"a", "b", "j", "b", "r", "a"}
	t.Log(StringsDedupWithSort(str1))
}

func TestStringsExcept(t *testing.T) {
	str1 := []string{"a", "b", "j", "b", "z", "a"}
	str2 := []string{"w", "b", "t", "b", "r", "a"}
	str3 := []string{"a", "e", "j", "b", "r", "a"}

	t.Log("1-2:", StringsExcept(str1, str2))
	t.Log("1-3:", StringsExcept(str1, str3))
	t.Log("3-2:", StringsExcept(str3, str2))

}
