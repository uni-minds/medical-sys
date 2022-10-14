/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: rand.go
 */

package tools

import (
	"math/rand"
	"time"
)

var r *rand.Rand

//const ALPHABET = "0123456789abcdef"
const ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//const alphacount = len(ALPHABET)

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func RandString0f(len int) string {
	return RandStringFromAlphabet(len, ALPHABET[:16])
}

func RandStringFromAlphabet(length int, alphabet string) string {
	alphaLen := 62
	if alphabet == "" {
		alphabet = ALPHABET
	} else {
		alphaLen = len(alphabet)
	}

	bs := make([]byte, length)
	for i := 0; i < length; i++ {
		b := r.Intn(alphaLen)
		bs[i] = alphabet[b]
	}
	return string(bs)
}
