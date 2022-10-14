/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: checksum_test.go
 */

package tools

import "testing"

func TestGetFileMD5(t *testing.T) {
	t.Log("M1")
	t.Log(GetFileMD5("checksum.go"))
}
