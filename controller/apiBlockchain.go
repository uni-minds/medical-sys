package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func GetBlockchainNodelist(ctx *gin.Context) {
	resp, _, _, _ := tools.HttpGet("http://localhost:10000/api/v1/node/list")

	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, tools.FailReturn(400, resp.Message))
	} else {
		ctx.JSON(http.StatusOK, tools.SuccessReturn(resp.Data))
	}
}

func GetBlockchainTPS(ctx *gin.Context) {
	ip := ctx.Query("ip")
	url := "http://localhost:10000/api/v1/node/tps"
	if ip != "" {
		url = fmt.Sprintf("http://%s:10000/api/v1/node/tps", ip)
	}
	resp, _, _, _ := tools.HttpGet(url)
	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, tools.FailReturn(400, resp.Message))
	} else {
		ctx.JSON(http.StatusOK, tools.SuccessReturn(resp.Data))
	}
}

func GetBlockHeight(ctx *gin.Context) {
	h := ctx.Query("height")
	url := fmt.Sprintf("http://localhost:10000/api/v1/block/record/height/%s", h)
	resp, _, _, _ := tools.HttpGet(url)
	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, tools.FailReturn(400, resp.Message))
	} else {
		ctx.JSON(http.StatusOK, tools.SuccessReturn(resp.Data))
	}
}
