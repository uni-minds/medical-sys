/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: common.go
 * Description:
 */

package tools

import (
	"fmt"
	"math/rand"
	"time"
	"uni-minds.com/liuxy/medical-sys/logger"
)

var r *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
const alphacount = len(alphabet)
const tag = "TOOL"

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(alphacount)
		bytes[i] = alphabet[b]
	}
	return string(bytes)
}

func log(level string, message ...interface{}) {
	if len(message) == 0 {
		logger.Write(tag, "i", level)
	} else {
		msg := ExpandInterface(message)
		logger.Write(tag, level, msg)
	}
}

func ExpandInterface(msg []interface{}) string {
	str := fmt.Sprintf("%v", msg[0])
	for _, v := range msg[1:] {
		str = fmt.Sprintf("%s %v", str, v)
	}
	return str
}
