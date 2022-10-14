/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: base.go
 */

package controller

import (
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/medical-sys/manager"
	"gitee.com/uni-minds/utils/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
	"time"
)

const tag = "CTRL"
const edaAddress = "localhost:80"       // http
const mboxGateway = "localhost:8442"    // https
const algoServer = "192.168.2.101:5000" // https

var randstr string
var randgen time.Time

func Init() (err error) {
	algofile = path.Join(global.GetPaths().Application, "algo.yaml")
	if err := tools.LoadYaml(algofile, &algolist); err != nil {
		algolist = global.DefaultAlgorithms()
		tools.SaveYaml(algofile, algolist)
	}

	var menuconfig = path.Join(global.GetPaths().Application, "menu.yaml")
	if err := tools.LoadYaml(menuconfig, &menudata); err != nil {
		menudata = global.DefaultMenuData()
		tools.SaveYaml(menuconfig, menudata)
	}

	return nil
}

func log(level string, message ...interface{}) {
	msg := tools.InterfaceExpand(message)
	logger.Write(tag, level, msg)
}

const (
	ETokenInvalid     = "登录凭证无效"
	EActionForbiden   = "禁止操作"
	EParameterInvalid = "参数异常"
)

func getRandomString() string {
	if time.Now().Sub(randgen) > 5*time.Second {
		randstr = tools.RandString0f(8)
		randgen = time.Now()
	}
	return randstr
}

func FailReturn(code int, msg interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	res["data"] = ""
	res["code"] = code
	res["msg"] = msg

	return res
}

// SuccessReturn api正确返回函数
func SuccessReturn(msg interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	res["data"] = msg
	res["code"] = http.StatusOK
	res["msg"] = "success"

	return res
}

func CookieWrite(ctx *gin.Context, key string, value string, age int) {
	cookieName := &http.Cookie{Name: key, Value: value, Path: "/", Secure: false, HttpOnly: true, MaxAge: age}
	http.SetCookie(ctx.Writer, cookieName)
}

func CookieRead(ctx *gin.Context, key string) (value string) {
	if cookie, err := ctx.Request.Cookie(key); err != nil {
		return ""
	} else {
		return cookie.Value
	}
}

func CookieRemove(ctx *gin.Context, key string) {
	clearCookieName := &http.Cookie{Name: key, Value: "", Path: "/", MaxAge: -1, Secure: false, HttpOnly: true}
	http.SetCookie(ctx.Writer, clearCookieName)
}

func CookieValidUid(ctx *gin.Context) (result bool, uid int) {
	uidstr := CookieRead(ctx, "uid")
	token := CookieRead(ctx, "token")

	if uidstr != "" && token != "" {
		uid, _ = strconv.Atoi(uidstr)
	}

	if manager.TokenValidator(uid, token) {
		return true, uid
	} else {
		return false, -1
	}
}
