/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: string.go
 */

package tools

import (
	_ "github.com/mattn/go-runewidth"
	"strings"
)

func LineBuilder(width int, char string) string {
	var bs strings.Builder

	for i := 0; i < width; i++ {
		bs.WriteString(char)
	}
	return bs.String()
}
