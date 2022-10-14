/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: db_medialock.go
 */

package manager

import (
	"errors"
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

const MediaAccessLockTime = 45 * time.Second

var mediaAccess MediaAccess

func MediaAccessSetLock(mediaIndex string, uid int, tp string) (status MediaLocker, err error) {
	if status, err = MediaAccessGetLock(mediaIndex); err != nil {
		// 未存在锁，或已超时
		mediaAccess.Lock.Lock()
		status = MediaLocker{
			Uid:  uid,
			Time: time.Now(),
			Type: tp,
		}
		mediaAccess.DB[mediaIndex] = status
		mediaAccess.Lock.Unlock()
		return status, nil

	} else {
		// 存在锁，未超时
		if status.Uid == uid {
			mediaAccess.Lock.Lock()
			status.Time = time.Now()
			status.Type = tp
			mediaAccess.DB[mediaIndex] = status
			mediaAccess.Lock.Unlock()
			return status, nil
		}
		return status, errors.New("lock by others")
	}
}
func MediaAccessGetLock(mediaIndex string) (status MediaLocker, err error) {
	mediaAccess.Lock.RLock()
	status, ok := mediaAccess.DB[mediaIndex]
	mediaAccess.Lock.RUnlock()
	if !ok || time.Now().Sub(status.Time) > MediaAccessLockTime {
		MediaAccessUnlock(mediaIndex, 0, true)
		err = errors.New("media lock not found")
	}
	return
}
func MediaAccessUnlock(mediaIndex string, uid int, override bool) bool {
	if override {
		mediaAccess.Lock.Lock()
		delete(mediaAccess.DB, mediaIndex)
		mediaAccess.Lock.Unlock()
		return true
	}

	mediaAccess.Lock.RLock()
	status, ok := mediaAccess.DB[mediaIndex]
	mediaAccess.Lock.RUnlock()
	if ok && status.Uid == uid {
		mediaAccess.Lock.Lock()
		delete(mediaAccess.DB, mediaIndex)
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
