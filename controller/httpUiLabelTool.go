/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpUiLabelTool.go
 */

/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  uiLabelTool
 * @Version: 1.0.0
 * @Date: 2020/4/10 09:18
 */

package controller

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/module"
	"gitee.com/uni-minds/medical-sys/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func UiLabeltoolGet(ctx *gin.Context) {
	mediaIndex := ctx.Param("mediaIndex")

	var summary module.MediaSummaryInfo
	if strings.Contains(mediaIndex, DICOM_TYPE_US_ID) {
		// instance_id
		summary, _ = module.InstanceGetSummary(mediaIndex)

	} else if len(mediaIndex) == 32 {
		mid := module.MediaGetMid(mediaIndex)
		summary, _ = module.MediaGetSummary(mid)
	}

	customScript := fmt.Sprintf("/webapp/labelsys/us/labelsys.js?rnd=%s", getRandomString())

	switch ctx.Param("usertype") {
	case "author", "review":
		ctx.HTML(http.StatusOK, "labelsys.html", gin.H{
			"title":          "影像标注 | Medi-sys",
			"mediaDuration":  summary.Duration,
			"mediaFrames":    summary.Frames,
			"mediaHeight":    summary.Height,
			"mediaWidth":     summary.Width,
			"mediaIndex":     mediaIndex,
			"custom_scripts": customScript,
		})
	}
}

var randstr string
var randgen time.Time

func getRandomString() string {
	if time.Now().Sub(randgen) > 5*time.Second {
		randstr = tools.RandString0f(8)
		randgen = time.Now()
	}
	return randstr
}
