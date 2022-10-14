package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"uni-minds.com/medical-sys/module"
)

func GroupGetHandler(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if !valid {
		ctx.JSON(http.StatusOK, FailReturn(ETokenInvalid))
		return
	}

	action := ctx.Query("action")
	switch action {
	case "getlist":
		gids := module.UserGetGroups(uid)
		for i, gid := range gids {
			if gid == 1 {
				gids = append(gids[0:i], gids[i+1:]...)
			}
		}

		if len(gids) == 0 {
			ctx.JSON(http.StatusOK, FailReturn("用户不属于任何组"))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(gids))
		}
	case "getname":
		gidstr := ctx.Query("gid")
		if gidstr == "" {
			ctx.JSON(http.StatusOK, FailReturn(EParameterInvalid))
			return
		}

		gid, err := strconv.Atoi(gidstr)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			return
		}

		dispname := module.GroupGetDisplayname(gid)
		ctx.JSON(http.StatusOK, SuccessReturn(dispname))
	default:
		ctx.JSON(http.StatusOK, FailReturn(action))
	}
}
