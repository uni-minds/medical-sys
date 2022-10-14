/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: apiSync.go
 * Date: 2022/09/07 18:41:07
 */

package controller

import (
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SyncPost(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}
	log("d", uid, "sync post")

	//module.PacsSync("192.168.3.101:8080")
	dh := database.BridgeGetPacsDatabaseHandler()
	store_group := ctx.Query("store_group")
	gid := module.GroupGetGid(store_group)
	if gid > 0 {
		stat := dh.Sync(false, 10, true)
		if len(stat.ImportStudies) > 0 {
			if err := module.GroupAddMedia(gid, stat.ImportStudies); err != nil {
				log("e", "Sync add dicom:", err.Error())
				ctx.JSON(http.StatusOK, FailReturn(403, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(stat))
			}
		}
	}
}
