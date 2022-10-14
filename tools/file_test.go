/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: file_test.go
 */

package tools

import "testing"

func TestCopyFile(t *testing.T) {
	CopyFile("base.go", "base-copy.go")

}
