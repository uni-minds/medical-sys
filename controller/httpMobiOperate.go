package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func MobiGetDevice(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "mobi_device.html", gin.H{
		"title":      "设备",
		"card_title": "算法推送",
		"page_id":    "Device",
	})
}

func MobiGetResult(ctx *gin.Context) {
	pipeline := ctx.Param("pipeline")
	url := fmt.Sprintf("http://localhost:10000/api/v1/mbox/pipeline/%s", pipeline)
	data, _, _, err := tools.HttpGet(url)

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
		})
	} else if data.Code != 200 {
		fmt.Println("B", data.Message)
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
		})
	}
}

func MobiMyExec(ctx *gin.Context) {
	dev := ctx.Query("dev")
	algoid := ctx.Query("algoid")
	url := fmt.Sprintf("https://localhost:8442/medi-box/%s/capture?remark=demo&algoid=%s", dev, algoid)

	data, _, _, _ := tools.HttpGet(url)
	pipeline := data.Data.(string)
	if pipeline == "" {
		ctx.JSON(http.StatusOK, FailReturn("pipeline is empty"))
		return
	}

	url = fmt.Sprintf("http://localhost:10000/api/v1/mbox/pipeline/%s", pipeline)
	for {
		data, _, _, err := tools.HttpGet(url)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			return
		} else if data.Data != nil {
			ctx.JSON(http.StatusOK, SuccessReturn(pipeline))
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func MobiRoot(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/mobi/device")
}
