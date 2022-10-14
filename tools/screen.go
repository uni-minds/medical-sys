/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: screen.go
 */

package tools

func ScreenClear() {
	print("\033[H\033[2J")
}
