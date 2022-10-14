/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  db_medialock
 * @Version: 1.0.0
 * @Date: 2020/4/9 13:23
 */

package manager

import (
	"errors"
	"log"
	"sync"
	"time"
)

type MediaLocker struct {
	Uid  int
	Time time.Time
	Type string
}
type MediaAccess struct {
	Lock sync.RWMutex
	DB   map[string]MediaLocker
}

var mediaAccessLockTime = 60 * time.Second
var mediaAccess MediaAccess

func mediaAccessLockInit() {
	mediaAccess.DB = make(map[string]MediaLocker, 0)
}
func MediaAccessSetLock(mediaHash string, uid int, tp string) (status MediaLocker, err error) {
	status, err = MediaAccessGetLock(mediaHash)
	if err == nil {
		// 存在锁，未超时
		if status.Uid == uid {
			mediaAccess.Lock.Lock()
			status.Time = time.Now()
			status.Type = tp
			mediaAccess.DB[mediaHash] = status
			mediaAccess.Lock.Unlock()
			log.Println("Locker time renew:", status)
		}
		return status, errors.New("media locked")
	} else {
		// 未存在锁，或已超时
		status = MediaLocker{
			Uid:  uid,
			Time: time.Now(),
			Type: tp,
		}
		mediaAccess.Lock.Lock()
		mediaAccess.DB[mediaHash] = status
		mediaAccess.Lock.Unlock()
		return status, nil
	}
}
func MediaAccessGetLock(mediaHash string) (status MediaLocker, err error) {
	mediaAccess.Lock.RLock()
	status, ok := mediaAccess.DB[mediaHash]
	mediaAccess.Lock.RUnlock()
	if !ok || time.Now().Sub(status.Time) > mediaAccessLockTime {
		MediaAccessUnlock(mediaHash, 0, true)
		err = errors.New("media lock not found")
	}
	return
}
func MediaAccessUnlock(mediaHash string, uid int, override bool) bool {
	if override {
		mediaAccess.Lock.Lock()
		delete(mediaAccess.DB, mediaHash)
		mediaAccess.Lock.Unlock()
		return true
	}

	mediaAccess.Lock.RLock()
	status, ok := mediaAccess.DB[mediaHash]
	mediaAccess.Lock.RUnlock()
	if ok && status.Uid == uid {
		mediaAccess.Lock.Lock()
		delete(mediaAccess.DB, mediaHash)
		mediaAccess.Lock.Unlock()
		return true
	}
	return false
}
func MediaAccessLockList() map[string]MediaLocker {
	mediaAccess.Lock.RLock()
	defer mediaAccess.Lock.RUnlock()
	return mediaAccess.DB
}
