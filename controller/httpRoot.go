package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/manager"
)

func RootGetHandler(ctx *gin.Context) {
	valid, _ := CookieValidUid(ctx)
	if valid {
		ctx.Redirect(http.StatusFound, "/ui/home")
	} else {
		ctx.Redirect(http.StatusFound, "/login")
	}
}

func RootUserLoginGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "userLogin.html", gin.H{
		"title":    "用户登录 ｜ Medi-sys",
		"loginapi": "/api/v1/login",
		"register": "/register",
	})
}

func RootUserLogoutGet(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if valid {
		manager.TokenRemove(uid)
		CookieRemove(ctx, "token")
	}
	ctx.Redirect(http.StatusFound, "/")
}
