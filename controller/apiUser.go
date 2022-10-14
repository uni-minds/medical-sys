package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/module"
)

type UserLastStatus struct {
	LastPageIndex int `json:"lastPageIndex"`
	LastGroupId   int `json:"lastGroupId"`
}

func UserGet(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if !valid {
		return
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
		ctx.JSON(http.StatusOK, FailReturn(action))
	}

}

func UserPost(ctx *gin.Context) {
	valid, _ := CookieValidUid(ctx)
	if !valid {
		ctx.JSON(http.StatusOK, FailReturn(ETokenInvalid))
		return
	}
}

func UserDelete(ctx *gin.Context) {
	valid, _ := CookieValidUid(ctx)
	if !valid {
		ctx.JSON(http.StatusOK, FailReturn(ETokenInvalid))
		return
	}
}

func UserPut(ctx *gin.Context) {

}
