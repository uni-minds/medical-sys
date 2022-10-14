/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiGroup.go
 */

package controller

import (
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GroupGet
//
// ?action=getlist [&type=*]
// ?action=getlistfull [&type=*]
// ?action=getname&gid=
func GroupGet(ctx *gin.Context) {
	uid := -1
	if uidi, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = uidi.(int)
	}

	action := ctx.Query("action")
	switch action {
	case "getlist":
		groupType := ctx.Query("type")
		gids := module.UserGetGroups(uid, groupType)

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

		groupType := ctx.Query("type")
		gids := module.UserGetGroups(uid, groupType)

		type ginfo struct {
			Gid   int
			Name  string
			GType string
		}

		data := make([]ginfo, 0)
		for _, gid := range gids {
			// remove administrators group
			if gid != 1 {
				data = append(data, ginfo{Gid: gid, Name: module.GroupGetDisplayname(gid), GType: module.GroupGetType(gid)})
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
