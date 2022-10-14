package manager

import (
	"sync"
	"time"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/tools"
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

func tokenInit() {
	tokenAccess.DB = make(map[int]TokenInfo, 0)
}

func TokenGenerater() string {
	return tools.GenSaltString(16, "0123456789abcdef")
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
