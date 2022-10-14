/*
 * Copyright (c) 2019-2021
 * Author: LIU Xiangyu
 * File: httpUiMediaScreen.go
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UiMediaScreenGet(ctx *gin.Context) {
	tp := ctx.Query("type")
	switch tp {
	case "us":
		switch ctx.Query("action") {
		case "screen":
			ctx.HTML(http.StatusOK, "mediascreen_us.html", gin.H{
				"title":          "超声挑图 | Medi-sys",
				"page_id":        "us-screen",
				"custom_scripts": "/webapp/js/us_pacs_screen_tool.js",
			})

		default:
			ctx.HTML(http.StatusOK, "mediascreen_us.html", gin.H{
				"title":          "超声挑图 | Medi-sys",
				"page_id":        "us-screen",
				"custom_scripts": "/webapp/js/us_pacs_screen_list.js",
			})
		}
	}
	return
}
