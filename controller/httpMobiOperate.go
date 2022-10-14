/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpMobiOperate.go
 */

package controller

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/utils/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type execInfo struct {
	Device  string `json:"dev"`
	AlgoRef string `json:"algo-ref"`
}

func MobiGetDevice(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "mobi_device.html", gin.H{
		"title":      "设备",
		"card_title": "算法推送",
		"page_id":    "Device",
	})
}

func MobiGetResult(ctx *gin.Context) {
	pipeline := ctx.Param("pipeline")
	url := fmt.Sprintf("http://%s/api/v1/mbox/pipeline/%s", edaAddress, pipeline)
	data, raw, err := tools.HttpGet(url)

	log("t", "url resp:", string(raw), err)
	if err != nil {
		ctx.HTML(http.StatusOK, "mobi_result.html", gin.H{
			"title":                  "结果",
			"page_id":                "Result",
			"card_title":             "消息详情",
			"pipeline":               pipeline,
			"tx_message_time":        "TMSG_T",
			"tx_message_content":     "TMSG_C",
			"tx_sandbox_id":          "SID",
			"tx_sandbox_type":        "ST",
			"tx_pipeline_init":       "PIPI",
			"tx_pipeline_upstream":   "PIPU",
			"tx_pipeline_downstream": "PIPD",
			"data_api":               url,
		})
	} else if data.Code != 200 {
		log("e", "mobiGetResult err code", data.Code)
	} else if data.Data == nil {
		ctx.HTML(http.StatusOK, "mobi_result.html", gin.H{
			"title":                  "结果",
			"page_id":                "Result",
			"card_title":             "未检索到相关消息",
			"pipeline":               pipeline,
			"tx_message_time":        "无",
			"tx_message_content":     "无法获取消息：消息不存在或正在共识队列中，请检查消息编号是否正确",
			"tx_sandbox_id":          "无",
			"tx_sandbox_type":        "none",
			"tx_pipeline_init":       "PIPI",
			"tx_pipeline_upstream":   "PIPU",
			"tx_pipeline_downstream": "PIPD",
			"data_api":               url,
		})
	} else {
		d := make([]map[string]string, 0)
		bs, _ := json.Marshal(data.Data)
		_ = json.Unmarshal(bs, &d)

		d0 := d[0]

		ctx.HTML(http.StatusOK, "mobi_result.html", gin.H{
			"title":                  "结果",
			"page_id":                "Result",
			"card_title":             "消息详情",
			"pipeline":               d0["Catalog"],
			"tx_message_time":        d0["Time"],
			"tx_message_content":     d0["Content"],
			"tx_sandbox_id":          d0["Upstream"],
			"tx_sandbox_type":        d0["Type"],
			"tx_pipeline_init":       d0["NodeId"],
			"tx_pipeline_upstream":   d0["SrcId"],
			"tx_pipeline_downstream": d0["DescId"],
			"data_api":               url,
		})
	}
}

func MobiMyExec(ctx *gin.Context) {
	var info execInfo
	err := ctx.BindJSON(&info)
	if err != nil {
		log("e", "E1", err.Error())
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return
	}

	info.AlgoRef = strings.Replace(info.AlgoRef, "//", "/", -1)

	url := fmt.Sprintf("https://%s/medi-box/%s/capture", mboxGateway, info.Device)
	log("i", "exec url", url)

	resp, _, err := tools.HttpPost(url, info, "json")
	if err != nil {
		log("e", "mbox-gateway no response:", err.Error())
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
	}

	pipeline := resp.Data.(string)
	if pipeline == "" {
		ctx.JSON(http.StatusOK, FailReturn(400, "pipeline is empty"))
		return
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(pipeline))
	}
}

func MobiRoot(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/mobi/device")
}
