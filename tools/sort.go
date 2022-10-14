/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  sort
 * @Version: 1.0.0
 * @Date: 2020/4/8 12:25
 */

package tools

import "fmt"

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
