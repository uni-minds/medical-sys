/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiAnalysis.go
 */

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/module"
)

func AnalysisCtPost(ctx *gin.Context) {
	class := ctx.Param("class")
	mode := ctx.Param("mode")
	switch class {
	case "ccta":
		switch mode {
		case "deepbuild":
			var data module.DeepBuild
			err := ctx.BindJSON(&data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else if aid, err := module.AiDemoCctaAnalysisDummy(data); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, "重建失败"))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(aid))
			}

		case "deepsearch":
			var data module.DeepSearch
			err := ctx.BindJSON(&data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else if sid, err := module.AiDemoCctaSearchDummy(data); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, "检索失败"))
			} else {
				//ctx.JSON(http.StatusOK, SuccessReturn(sid))
				ctx.JSON(http.StatusOK, FailReturn(500, fmt.Sprintf("未检测到匹配项，请扩充特征池：%s", sid)))
			}
		}
	}
}
