/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: string.go
 */

package tools

import (
	"encoding/json"
	_ "github.com/mattn/go-runewidth"
	"sort"
	"strings"
)

func LineBuilder(width int, char string) string {
	var bs strings.Builder

	for i := 0; i < width; i++ {
		bs.WriteString(char)
	}
	return bs.String()
}

func StringCompress(strs []string) (str string, err error) {
	bs, err := json.Marshal(strs)
	return string(bs), err
}

func StringDecompress(str string) (strs []string, err error) {
	if str == "" {
		return make([]string, 0), nil
	}

	bs := []byte(str)
	err = json.Unmarshal(bs, &strs)
	return strs, err
}

func StringsDedup(slice []string) []string {
	i := 0

	var j int
	for {
		if i >= len(slice)-1 {
			break
		}

		for j = i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				slice = append(slice[:j], slice[j+1:]...)
			}
		}
		i++
	}
	return slice
}

func StringsDedupWithSort(a []string) (ret []string) {
	sort.Strings(a)
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func StringsExcept(a []string, b []string) (ret []string) {
	strx := StringsDedupWithSort(a)
	stry := StringsDedupWithSort(b)

	iy := 0
	ix := 0
	for ix < len(strx) {

		if iy >= len(stry) {
			ret = append(ret, strx[ix:]...)
			break
		}

		if strx[ix] < stry[iy] {
			ret = append(ret, strx[ix])
			ix++
		} else if strx[ix] == stry[iy] {
			ix++
			iy++
		} else if strx[ix] > stry[iy] {
			iy++
		}

	}
	return ret
}
