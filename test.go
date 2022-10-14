package main

import (
	"encoding/json"
	"fmt"
)

/*
type MapSorter []MapItem

func NewMapSorter(m map[string]string) MapSorter {
	ms := make(MapSorter,0,len(m))
	for k,v:=range m {
		ms=append(ms,MapItem{Key:k,Value:v})
	}
	return ms
}

type MapItem struct {
	Key string
	Value string
}

func (ms MapSorter) Len() int{
	return len(ms)
}

func (ms MapSorter) Swap(i,j int) {
	ms[i],ms[j]=ms[j],ms[i]
}

func (ms MapSorter) Less(i,j int) bool {
	return ms[i].Key <ms[j].Key
}

*/

func main() {
	a := "[1,2,3,4]"
	b := "[5,6,7,8]"
	type abc struct {
		A []int
		B []int
	}
	var tmp []int
	var tmp2 abc
	json.Unmarshal([]byte(a), &tmp2.A)
	//tmp2.A=tmp
	fmt.Println(tmp, tmp2)
	tmp = []int{}
	json.Unmarshal([]byte(b), &tmp2.B)
	//tmp2.B=tmp
	fmt.Println(tmp, tmp2)

}
