package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

var algolist []global.AlgorithmInfo
var algofile string

func init() {
	algofile = path.Join(global.GetAppSettings().SystemAppPath, "algo.yaml")
	if err := tools.LoadYaml(algofile, &algolist); err != nil {
		algolist = global.DefaultAlgorithms()
		tools.SaveYaml(algofile, algolist)
	}
}

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
