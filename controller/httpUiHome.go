/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpUiHome.go
 */

package controller

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UiHomeGet(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	bgContent := fmt.Sprintf("您好，%s老师。请单击左侧菜单栏选择相应的功能。", module.UserGetRealname(uid))

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "首页 ｜ Medi-Sys",
		"page_id":    "index",
		"bg_content": bgContent,
	})
}

func UiMediaListGet(ctx *gin.Context) {
	tp := ctx.Query("type")
	switch tp {
	case "us":
		customScript := fmt.Sprintf("/webapp/medialist/us/medialist.js?rnd=%s", getRandomString())
		ctx.HTML(http.StatusOK, "medialist_us.html", gin.H{
			"title":          "超声影像检索 | Medi-sys",
			"page_id":        "us-medialist",
			"custom_scripts": customScript,
		})
		break

	case "ct":
		customScript := fmt.Sprintf("/webapp/medialist/ct/medialist.js?rnd=%s", getRandomString())
		ctx.HTML(http.StatusOK, "medialist_ct.html", gin.H{
			"title":          "CT影像检索 | Medi-sys",
			"page_id":        "ct-medialist",
			"custom_scripts": customScript,
		})
		break
	}
	return
}

func UiMediaScreenGet(ctx *gin.Context) {
	tp := ctx.Query("type")
	switch tp {
	case "us":
		switch ctx.Query("action") {
		case "screen":
			customScript := fmt.Sprintf("/webapp/medialist/us/us_pacs_screen_tool.js?rnd=%s", getRandomString())
			ctx.HTML(http.StatusOK, "mediascreen_us.html", gin.H{
				"title":          "超声挑图 | Medi-sys",
				"page_id":        "us-screen",
				"custom_scripts": customScript,
			})

		default:
			customScript := fmt.Sprintf("/webapp/medialist/us/us_pacs_screen_list.js?rnd=%s", getRandomString())
			ctx.HTML(http.StatusOK, "mediascreen_us.html", gin.H{
				"title":          "超声挑图 | Medi-sys",
				"page_id":        "us-screen",
				"custom_scripts": customScript,
			})
		}
	}
	return
}
