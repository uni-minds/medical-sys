/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: db_token.go
 */

package manager

import (
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/tools"
	"sync"
	"time"
)

type TokenInfo struct {
	Token string
	Time  time.Time
}

type TokenAccess struct {
	Lock sync.RWMutex
	DB   map[int]TokenInfo
}

var tokenAccess TokenAccess

func TokenGenerater() string {
	return tools.RandString0f(16)
}

func TokenNew(uid int) (token string) {
	token = TokenGenerater()
	tokenAccess.Lock.Lock()
	tokenAccess.DB[uid] = TokenInfo{
		Token: token,
		Time:  time.Now(),
	}
	tokenAccess.Lock.Unlock()
	database.UserSetToken(uid, token)
	return token
}

func TokenRemove(uid int) {
	tokenAccess.Lock.Lock()
	delete(tokenAccess.DB, uid)
	tokenAccess.Lock.Unlock()
	database.UserSetToken(uid, "")
}

func TokenValidator(uid int, token string) bool {
	if uid < 0 || token == "" {
		return false
	}

	tokenAccess.Lock.RLock()
	defer tokenAccess.Lock.RUnlock()
	data, ok := tokenAccess.DB[uid]
	if ok {
		return token == data.Token

	} else {
		r := database.UserTokenCheck(uid, token)
		if r {
			tokenAccess.DB[uid] = TokenInfo{
				Token: token,
				Time:  time.Now(),
			}
		}
		return r
	}
}

func TokenList() map[int]TokenInfo {
	tokenAccess.Lock.RLock()
	defer tokenAccess.Lock.RUnlock()
	return tokenAccess.DB
}
