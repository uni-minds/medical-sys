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
	"uni-minds.com/liuxy/medical-sys/module"
)

func UILabeltoolGet(ctx *gin.Context) {
	content := ""
	tp := ctx.Query("type")
	switch tp {
	case "us":

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

		ctx.HTML(http.StatusOK, "labelsys-v2.html", gin.H{
			"title":          "影像标注 | Medi-sys",
			"media_hash":     mediaHash,
			"media_duration": summary.Duration,
			"media_frames":   summary.Frames,
			"media_height":   summary.Height,
			"media_width":    summary.Width,
			"custom_scripts": "/webapp/us/js/labelsys-v2.js",
		})
		break

	default:
		content = fmt.Sprintf("未知类型：%s", tp)
	}

	ctx.Writer.Write([]byte(content))
	return
}
