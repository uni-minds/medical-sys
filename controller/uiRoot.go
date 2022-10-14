package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/manager"
	"uni-minds.com/liuxy/medical-sys/module"
)

func UIRootGetHandler(ctx *gin.Context) {
	valid, _ := CookieValidUid(ctx)
	if valid {
		ctx.Redirect(http.StatusFound, "/ui/home")
	} else {
		ctx.HTML(http.StatusOK, "userLogin.html", gin.H{
			"title":    "用户登录 ｜ Medi-sys",
			"loginapi": "/api/login",
			"register": "/ui/register",
		})
	}
}
func UIHomeGetHandler(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if valid {
		bg_content := fmt.Sprintf("您好，%s。\n"+
			"单击左侧菜单栏选择相应的功能。", module.UserGetRealname(uid))

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "首页 ｜ Medi-Sys",
			"page_id":    "index",
			"bg_content": bg_content,
		})
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}
func UIRegisterGetHandler(ctx *gin.Context) {
	if global.GetAppSettings().EnableUserRegister {
		ctx.HTML(http.StatusOK, "userRegister.html", gin.H{
			"title":  "Medi-sys | 用户注册",
			"regapi": "/api/register",
		})
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}
func UILogoutGetHandler(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if valid {
		manager.TokenRemove(uid)
		CookieRemove(ctx, "token")
	}
	ctx.Redirect(http.StatusFound, "/")
}
