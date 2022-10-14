package manager

import (
	"gitee.com/uni-minds/medical-sys/logger"
)

var log *logger.Logger

func Init() (err error) {
	log = logger.NewLogger("TOKN")
	log.Println("init: token manager")
	tokenAccess.DB = make(map[int]TokenInfo, 0)
	mediaAccess.DB = make(map[string]MediaLocker, 0)
	return nil
}
