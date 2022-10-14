/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiAI.go
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/module"
)

type AlgoParamsData struct {
	Params string
}

func AiAlgoGet(ctx *gin.Context) {
	// :modal/:class/:algo/:aid/:part
	modal := ctx.Param("modal")
	class := ctx.Param("class")
	algo := ctx.Param("algo")
	aid := ctx.Param("aid")
	part := ctx.Param("part")
	log("t", "ai get:", modal, class, algo, aid, part)

	switch modal {
	case "ct":
		switch class {
		case "ccta":
			result := module.AlgoCctaGetFeatureResult(aid, part)
			if result != nil {
				ctx.JSON(http.StatusOK, SuccessReturn(result))
			} else {
				ctx.JSON(http.StatusOK, FailReturn(400, "CCTA Result failed"))
			}
		}
	}
}

func AiAlgoPost(ctx *gin.Context) {
	modal := ctx.Param("modal")
	class := ctx.Param("class")
	algo := ctx.Param("algo")

	var data AlgoParamsData
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return
	}

	var aid string
	log("t", "ai post:", modal, class, algo)
	switch modal {
	case "ct":
		switch class {
		case "cta":
			aid, err = module.RunAlgo(algo, data.Params)
		case "ccta":
			aid, err = module.RunAlgo(algo, data.Params)
		}

	case "us":

	}
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(aid))
	}
}
