/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  sort_test.go
 * @Version: 1.0.0
 * @Date: 2020/4/8 12:29
 */

package tools

import (
	"fmt"
	"sort"
	"testing"
)

func TestStringSorterInt(t *testing.T) {
	data := make(map[int]interface{}, 0)
	data[3] = -1
	data[2] = 2
	data[1] = 2
	data[4] = 2
	data[100] = 4

	m := MediaSorter(data)
	sort.Sort(m)
	fmt.Println(m)
	//sort.Sort(sort.Reverse(m))
	//fmt.Println(m)
}
