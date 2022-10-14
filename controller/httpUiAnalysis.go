/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpUiAnalysis.go
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UiAnalysisGet(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	switch q.Get("type") {
	case "deepbuild":
		switch q.Get("mode") {
		case "cta":
			ctx.HTML(http.StatusOK, "analysis_cta.html", nil)
			return
		case "ccta":
			ctx.HTML(http.StatusOK, "analysis_ccta.html", nil)
			return
		case "ccta3d":
			ctx.HTML(http.StatusOK, "analysis_ccta_3d.html", nil)
			return

		}

	case "deepsearch":
		switch q.Get("mode") {
		case "ccta":
			ctx.HTML(http.StatusOK, "analysis_deepsearch.html", nil)
			return
		case "cta":
			ctx.HTML(http.StatusOK, "analysis_deepsearch_cta.html", nil)
			return
		}
	}
}
