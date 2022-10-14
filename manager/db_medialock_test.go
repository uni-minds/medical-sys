/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  db_medialock_test.go
 * @Version: 1.0.0
 * @Date: 2020/4/9 21:51
 */

package manager

import (
	"testing"
	"time"
)

func TestMediaAccessGetLock(t *testing.T) {
	mediaAccessLockInit()
	status, err := MediaAccessSetLock("AA", 1, "AU")
	t.Log(status, err, MediaAccessLockList())
	time.Sleep(4 * time.Second)
	status, err = MediaAccessSetLock("AA", 2, "BU")
	t.Log(status, err, MediaAccessLockList())
	time.Sleep(4 * time.Second)
	status, err = MediaAccessGetLock("AA")
	t.Log(status, err, MediaAccessLockList())
}
