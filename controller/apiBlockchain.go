/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiBlockchain.go
 */

package controller

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetBlockchainNodelist(ctx *gin.Context) {
	resp, _, _ := tools.HttpGet(fmt.Sprintf("http://%s/api/v1/node/list", edaAddress))

	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, FailReturn(400, resp.Message))
	} else {
		var list []global.NodeInfo
		bs, _ := json.Marshal(resp.Data)
		err := json.Unmarshal(bs, &list)
		if err != nil {
			log("e", "get nodelist:", err.Error())
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		} else {
			var callback []global.NodeStatus
			for _, info := range list {
				alive := time.Now().Sub(info.LastTalk) < 35*time.Second
				if !alive {
					log("w", "node timeout:", info.Name, time.Now().Sub(info.LastTalk).Seconds(), info.LastTalk)
				}
				s := global.NodeStatus{
					Name:   info.Name,
					Alive:  alive,
					IP:     info.IP,
					Height: info.BlockHeight,
				}
				if info.IP == "localhost" {
					callback = append([]global.NodeStatus{s}, callback...)
				} else {
					callback = append(callback, s)
				}
			}
			ctx.JSON(http.StatusOK, SuccessReturn(callback))
		}
	}
}

func GetBlockchainTPS(ctx *gin.Context) {
	ip := ctx.Query("addr")
	if ip == "" {
		ip = edaAddress
	}
	url := fmt.Sprintf("http://%s/api/v1/node/tps", ip)
	resp, _, _ := tools.HttpGet(url)
	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, FailReturn(400, resp.Message))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(resp.Data))
	}
}

func GetBlockHeight(ctx *gin.Context) {
	h := ctx.Query("height")
	resp, _, _ := tools.HttpGet(fmt.Sprintf("http://%s/api/v1/block/record/height/%s", edaAddress, h))
	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, FailReturn(400, resp.Message))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(resp.Data))
	}
}
