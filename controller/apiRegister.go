/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiRegister.go
 */

package controller

import (
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/module"
)

const (
	Ereg00CommonError     = "异常（Erx00）"
	Ereg01InvalidData     = "无效信息（Erx01）"
	Ereg02InvalidRegcode  = "无效邀请码（Erx02）"
	Ereg03UsernameExisted = "用户名已存在（Erx03）"
)

type clientPushRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Realname string `json:"realname"`
	Regcode  string `json:"regcode"`
}

func RootUserRegisterGet(ctx *gin.Context) {
	if global.GetAppSettings().UserRegisterEnable {
		ctx.HTML(http.StatusOK, "userRegister.html", gin.H{
			"title":  "用户注册 | Medi-sys",
			"regapi": "/register",
		})
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func RootUserRegisterPost(ctx *gin.Context) {
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
