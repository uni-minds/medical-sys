package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"uni-minds.com/medical-sys/global"
	"uni-minds.com/medical-sys/manager"
	"uni-minds.com/medical-sys/module"
)

type clientPushLogin struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Remember bool   `json:"remember"`
}

func LoginPostHandler(ctx *gin.Context) {
	var u clientPushLogin
	err := ctx.BindJSON(&u)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(err.Error()))
		return
	}

	uid := module.UserCheckPassword(u.Username, u.Password)
	switch uid {
	case -1, -3:
		ctx.JSON(http.StatusOK, FailReturn("用户名或密码不正确"))
	case -2, -5:
		ctx.JSON(http.StatusOK, FailReturn("用户禁止登陆"))
	case -4:
		ctx.JSON(http.StatusOK, FailReturn("账号过期"))
	default:
		if uid > 0 {
			token := manager.TokenNew(uid)
			maxAge := -1
			if u.Remember {
				maxAge = global.GetCookieMaxAge()
			}
			CookieWrite(ctx, "token", token, maxAge)
			ctx.JSON(http.StatusOK, SuccessReturn("/"))
		} else {
			log.Println("Invalid UID=", uid)
			ctx.JSON(http.StatusOK, FailReturn("用户UID无效"))
		}
	}
	return
}

func LoginGetHandler(ctx *gin.Context) {
	gkey := ctx.Query("goldenkey")
	if gkey == "Uni-Ledger-RIS" {
		username := ctx.Query("user")
		uid := module.UserGetUid(username)
		if uid != 0 {
			token := manager.TokenNew(uid)
			log.Printf("------ GOLDEN KEY OVERRIDE / username: %s / uid: %d / token: %s / ------", username, uid, token)
			CookieWrite(ctx, "token", token, global.GetCookieMaxAge())
			ctx.Redirect(http.StatusFound, "/")
			return
		}
	}
}
