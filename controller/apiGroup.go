/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiGroup.go
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"uni-minds.com/liuxy/medical-sys/module"
)

func GroupGet(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if !valid {
		ctx.JSON(http.StatusOK, FailReturn(400, ETokenInvalid))
		return
	}

	action := ctx.Query("action")
	switch action {
	case "getlist":
		gids := module.UserGetGroups(uid, ctx.Query("grouptype"))
		for i, gid := range gids {
			// remove administrators group
			if gid == 1 {
				gids = append(gids[0:i], gids[i+1:]...)
			}
		}

		if len(gids) == 0 {
			ctx.JSON(http.StatusOK, FailReturn(400, "用户不属于任何组"))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(gids))
		}

	case "getlistfull":
		gids := module.UserGetGroups(uid, ctx.Query("grouptype"))

		type ginfo struct {
			Gid  int
			Name string
		}
		data := make([]ginfo, 0)
		for _, gid := range gids {
			// remove administrators group
			if gid != 1 {
				data = append(data, ginfo{Gid: gid, Name: module.GroupGetDisplayname(gid)})
			}
		}

		if len(gids) == 0 {
			ctx.JSON(http.StatusOK, FailReturn(400, "用户不属于任何组"))
			return
		}

		ctx.JSON(http.StatusOK, SuccessReturn(data))

	case "getname":
		gidstr := ctx.Query("gid")
		if gidstr == "" {
			ctx.JSON(http.StatusOK, FailReturn(400, EParameterInvalid))
			return
		}

		gid, err := strconv.Atoi(gidstr)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			return
		}

		dispname := module.GroupGetDisplayname(gid)
		ctx.JSON(http.StatusOK, SuccessReturn(dispname))

	default:
		ctx.JSON(http.StatusOK, FailReturn(400, action))
	}
}
