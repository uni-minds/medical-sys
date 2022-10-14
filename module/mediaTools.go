/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  mediaTools
 * @Version: 1.0.0
 * @Date: 2020/3/31 01:01
 */

package module

import (
	"encoding/json"
	"sort"
	"strconv"
)

func parseUidLidStringToDbMap(juidstr, jlidstr string) map[int]int {
	var uidstrs, lidstrs []string
	db := make(map[int]int, 0)

	_ = json.Unmarshal([]byte(juidstr), &uidstrs)
	_ = json.Unmarshal([]byte(jlidstr), &lidstrs)

	for i, uidstr := range uidstrs {
		uid, _ := strconv.Atoi(uidstr)
		lid, _ := strconv.Atoi(lidstrs[i])
		db[uid] = lid
	}
	return db
}

func parseUidLidMapMapToString(db map[int]int) (juidstr, jlidstr string) {
	var uids []int
	var uidstrs, lidstrs []string
	for key, _ := range db {
		uids = append(uids, key)
	}
	sort.Ints(uids)
	for _, uid := range uids {
		uidstrs = append(uidstrs, strconv.Itoa(uid))
		lidstrs = append(lidstrs, strconv.Itoa(db[uid]))
	}
	jbUids, _ := json.Marshal(uidstrs)
	jbLids, _ := json.Marshal(lidstrs)
	return string(jbUids), string(jbLids)
}
