/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpRoot.go
 */

package controller

import (
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/manager"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Ereg00CommonError     = "异常（Erx00）"
	Ereg01InvalidData     = "无效信息（Erx01）"
	Ereg02InvalidRegcode  = "无效邀请码（Erx02）"
	Ereg03UsernameExisted = "用户名已存在（Erx03）"
)

func RootGet(ctx *gin.Context) {
	switch ctx.Param("op") {
	case "login":
		ctx.HTML(http.StatusOK, "userLogin.html", gin.H{
			"title":     "用户登录 ｜ Medi-sys",
			"masterApi": "/api/user/login",
		})

	case "logout":
		valid, uid := CookieValidUid(ctx)
		if valid {
			manager.TokenRemove(uid)
			CookieRemove(ctx, "token")
		}
		ctx.Redirect(http.StatusFound, "/")

	case "forget":
		ctx.HTML(http.StatusOK, "userForget.html", gin.H{
			"title":     "重置密码 | Medi-sys",
			"masterApi": "/api/user/forget",
		})

	case "register":
		if !global.GetAppSettings().UserRegisterEnable {
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		ctx.HTML(http.StatusOK, "userRegister.html", gin.H{
			"title":     "用户注册 | Medi-sys",
			"masterApi": "/api/user/register",
		})

	default:
		valid, _ := CookieValidUid(ctx)
		if valid {
			ctx.Redirect(http.StatusFound, "/ui/home")
		} else {
			ctx.Redirect(http.StatusFound, "/login")
		}
	}
}
