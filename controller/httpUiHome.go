/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpUiHome.go
 */

package controller

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UiHomeGet(ctx *gin.Context) {
	js := "/dist/js/home.js"
	if global.FlagGetDebug() {
		js = fmt.Sprintf("%s?rnd=%s", js, getRandomString())
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":          "首页 ｜ Medi-Sys",
		"custom_scripts": js,
	})
}

func UiScreenSeriesGet(ctx *gin.Context) {
	js := "/dist/js/pacs_us_tool.js"
	if global.FlagGetDebug() {
		js = fmt.Sprintf("%s?rnd=%s", js, getRandomString())
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":          "病例筛查 | Medi-sys",
		"custom_scripts": js,
	})
}

func UiLabeltoolGet(ctx *gin.Context) {
	js := "/dist/js/labelsys_core.js"
	if global.FlagGetDebug() {
		js = fmt.Sprintf("%s?rnd=%s", js, getRandomString())
	}

	switch ctx.Param("usertype") {
	case "author", "review":
		ctx.HTML(http.StatusOK, "labelsys.html", gin.H{
			"custom_scripts": js,
		})
	}
}
