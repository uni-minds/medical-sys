package tools

import "testing"

func TestStringCompress(t *testing.T) {
	//data1 := []string{"a","b","c"}
	data1 := []string{}
	str, err := StringCompress(data1)
	t.Log(str, err)

	data2 := ""
	strD, err := StringDecompress(data2)
	t.Log(strD, err)
}
