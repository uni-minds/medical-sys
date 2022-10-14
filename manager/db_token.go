package manager

import (
	"sync"
	"time"
	"uni-minds.com/medical-sys/tools"
)

type TokenInfo struct {
	Token string
	Uid   int
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
	return tools.GenSaltString(16)
}

func TokenNew(uid int) (token string) {
	token = TokenGenerater()
	tokenAccess.Lock.Lock()
	tokenAccess.DB[uid] = TokenInfo{
		Token: token,
		Uid:   uid,
		Time:  time.Now(),
	}
	tokenAccess.Lock.Unlock()
	return token
}

func TokenRemove(uid int) {
	tokenAccess.Lock.Lock()
	delete(tokenAccess.DB, uid)
	tokenAccess.Lock.Unlock()
}

func TokenValidator(token string) (uid int) {
	tokenAccess.Lock.RLock()
	defer tokenAccess.Lock.RUnlock()
	for _, v := range tokenAccess.DB {
		if token == v.Token {
			return v.Uid
		}
	}
	return
}

func TokenList() map[int]TokenInfo {
	tokenAccess.Lock.RLock()
	defer tokenAccess.Lock.RUnlock()
	return tokenAccess.DB
}
