/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiRawData.go
 */

package controller

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"github.com/gin-gonic/gin"
	"net/http"
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
