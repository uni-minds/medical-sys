/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  mediaTools_test.go
 * @Version: 1.0.0
 * @Date: 2020/3/31 01:01
 */

package module

import (
	"testing"
)

func Test_parseUidLidStringToDbMap(t *testing.T) {
	uidstr := `["2","3","1"]`
	lidstr := `["20","30","10"]`
	db := parseUidLidStringToDbMap(uidstr, lidstr)
	t.Log(db)
	db[4] = 40
	db[0] = 0
	sa, sb := parseUidLidMapMapToString(db)
	t.Log(sa, sb)
}
