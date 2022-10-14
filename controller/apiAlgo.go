/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiAlgo.go
 */

package controller

import (
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

var algolist []global.AlgorithmInfo
var algofile string

func AlgoGet(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, SuccessReturn(algolist))
}

func AlgoPost(ctx *gin.Context) {
	var algo global.AlgorithmInfo
	err := ctx.BindJSON(&algo)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
	} else {
		for _, val := range algolist {
			if val.Name == algo.Name {
				ctx.JSON(http.StatusOK, FailReturn(400, "same name exist"))
				return
			}
		}

		algolist = append(algolist, global.AlgorithmInfo{
			Index: len(algolist) + 1,
			Name:  algo.Name,
			Ref:   algo.Ref,
		})
		tools.SaveYaml(algofile, algolist)
		ctx.JSON(http.StatusOK, SuccessReturn(algo.Name))
	}
}
