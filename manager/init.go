package manager

import "fmt"

func init() {
	fmt.Println("module init: manager")
	tokenAccess.DB = make(map[int]TokenInfo, 0)
	mediaAccess.DB = make(map[string]MediaLocker, 0)
}
