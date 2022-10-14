/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: sort.go
 */

package tools

import (
	"fmt"
	"sort"
)

type mapSorter []MapItem
type MapItem struct {
	Mid   int
	Value interface{}
}

func MediaSorter(m map[int]interface{}) mapSorter {
	ms := make(mapSorter, 0, len(m))
	for key, value := range m {
		ms = append(ms, MapItem{
			Mid:   key,
			Value: value,
		})
	}
	return ms
}

func (ms mapSorter) Len() int {
	return len(ms)
}

func (ms mapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms mapSorter) Less(i, j int) bool {
	if ms[i].Value == ms[j].Value {
		return ms[i].Mid < ms[j].Mid
	} else {
		switch ms[i].Value.(type) {
		case string:
			return ms[i].Value.(string) < ms[j].Value.(string)
		case int:
			return ms[i].Value.(int) < ms[j].Value.(int)
		case float64:
			return ms[i].Value.(float64) < ms[j].Value.(float64)
		default:
			fmt.Println("Unknow sort type", ms[i].Value)
			return false
		}
	}
}

func RemoveDuplicateInt(a []int) []int {
	sort.Ints(a)
	i := 0
	for j := 1; j < len(a); j++ {
		if a[i] != a[j] {
			i++
			a[i] = a[j]
		}
	}
	return a[:i+1]
}

func RemoveElementInt(a []int, ele int) []int {
	a = RemoveDuplicateInt(a)
	for k, v := range a {
		if ele == v {
			return append(a[:k], a[k+1:]...)
		}
	}
	return a
}
