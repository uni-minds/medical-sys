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
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"uni-minds.com/liuxy/medical-sys/module"
)

func UILabeltoolGetHandler(ctx *gin.Context) {
	valid, _ := CookieValidUid(ctx)
	if !valid {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	content := ""
	tp := ctx.Query("type")
	switch tp {
	case "us":
		crfScript := ""
		switch strings.ToLower(ctx.Query("crf")) {
		case "4ap":
			crfScript = "/webapp/ultrasonic/js/crf_4ap.js"
		case "a":
			crfScript = "/webapp/ultrasonic/js/crf_a.js"
		case "l":
			crfScript = "/webapp/ultrasonic/js/crf_l.js"
		case "r":
			crfScript = "/webapp/ultrasonic/js/crf_r.js"
		default:
			crfScript = "/webapp/ultrasonic/js/crf_4ap.js"
		}

		mediaHash := ctx.Query("media")
		if mediaHash == "" {
			return
		}
		mid := module.MediaGetMid(mediaHash)

		switch ctx.Query("action") {
		case "author":

		case "review":

		}
		//
		summary, _ := module.MediaGetSummary(mid)

		ctx.HTML(http.StatusOK, "ultrasonic_labelsys-v1.html", gin.H{
			"title":          "影像标注 | Medi-sys",
			"media_hash":     mediaHash,
			"media_duration": summary.Duration,
			"media_frames":   summary.Frames,
			"media_height":   summary.Height,
			"media_width":    summary.Width,
			"label_data":     module.LabelGetJson(mediaHash),
			"crf_scripts":    crfScript,
			"custom_scripts": "/webapp/ultrasonic/js/ultrasonic_labelsys-v1-mod-forV2.js",
		})
		break

	default:
		content = fmt.Sprintf("未知类型：%s", tp)
	}

	ctx.Writer.Write([]byte(content))
	return
}
