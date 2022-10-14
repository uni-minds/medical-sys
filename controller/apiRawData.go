/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiRawData.go
 */

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

var menudata []global.MenuStruct

func RawDataGet(ctx *gin.Context) {
	action := ctx.Query("action")
	switch action {
	// raw?action=getversion
	case "getversion":
		ctx.JSON(http.StatusOK, SuccessReturn(global.GetCopyrightHtml()))
		return

	// raw?action=getmenujson
	case "getmenujson":
		ctx.JSON(http.StatusOK, SuccessReturn(getMenuData()))
		return

	// raw?action=getviewjson&view=4ap
	case "getviewjson":
		view := ctx.Query("view")
		if view != "" {
			ctx.JSON(http.StatusOK, SuccessReturn(global.DefaultUltrasonicViewData(view)))
		} else {
			ctx.JSON(http.StatusOK, FailReturn(400, "N/A"))
		}
		return

	default:
		ctx.JSON(http.StatusOK, FailReturn(400, fmt.Sprintf("unknown action: %s", action)))
		return
	}
}

func getMenuData() []global.MenuStruct {
	return menudata
}

func init() {
	var menuconfig = path.Join(global.GetAppSettings().SystemAppPath, "menu.yaml")
	if err := tools.LoadYaml(menuconfig, &menudata); err != nil {
		menudata = global.DefaultMenuData()
		tools.SaveYaml(menuconfig, menudata)
	}
}
