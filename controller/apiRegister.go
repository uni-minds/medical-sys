package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"uni-minds.com/medical-sys/global"
	"uni-minds.com/medical-sys/module"
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

func RegisterPost(ctx *gin.Context) {
	var r clientPushRegister
	err := ctx.BindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(Ereg01InvalidData))
		return
	}

	log.Println("注册用户", r)
	if r.Regcode != global.GetUserRegCode() {
		ctx.JSON(http.StatusOK, FailReturn(Ereg02InvalidRegcode))
		return
	}

	uid := module.UserGetUid(r.Username)
	if uid != 0 {
		ctx.JSON(http.StatusOK, FailReturn(Ereg03UsernameExisted))
		return
	}

	err = module.UserCreate(r.Username, r.Password, r.Email, r.Realname, "REG")
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(Ereg00CommonError))
	} else {
		module.UserSetActive(module.UserGetUid(r.Username))
		ctx.JSON(http.StatusOK, SuccessReturn("/"))
	}
}
