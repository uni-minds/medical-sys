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
	"log"
	"net/http"
	"uni-minds.com/medical-sys/module"
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
		switch ctx.Query("crf") {
		case "4ap":
			crfScript = "/webapp/ultrasonic/js/crf_4ap.js"
		case "A":
			crfScript = "/webapp/ultrasonic/js/crf_a.js"
		default:
			crfScript = "/webapp/ultrasonic/js/crf_4ap.js"
			ctx.Writer.Write([]byte(content))
			return
		}

		mediaHash := ctx.Query("media")
		labelUUID := ctx.Query("label")
		if mediaHash == "" {
			return
		}
		mid := module.MediaGetMid(mediaHash)

		authorJson := "{}"
		reviewJson := "{}"
		switch ctx.Query("action") {
		case "author":
			if len(labelUUID) != 32 {
				// 创建 LabelUUID

			} else {
				// 编辑或预览
				authorJson = module.LabelGetAuthorJson(labelUUID)
			}
		case "review":
			if len(labelUUID) != 32 {
				// 新的审阅
				summary, err := module.MediaGetSummary(mid)
				if err != nil {
					log.Println(err.Error())
				}
				authorJson = module.LabelGetAuthorJson(summary.AuthorLids[0])

			} else {
				// 修改审阅
				authorJson, reviewJson, _ = module.LabelGetReviewJson(labelUUID)
			}
		}
		//
		summary, _ := module.MediaGetSummary(mid)
		//authorData, _ := module.LabelGetAuthorJson(mid, -1)
		//reviewerData, _ := module.LabelGetReviewJson(mid, uid)

		ctx.HTML(http.StatusOK, "ultrasonic_labelsys-v1.html", gin.H{
			"title":          "影像标注 | Medi-sys",
			"media_hash":     mediaHash,
			"media_duration": summary.Duration,
			"media_frames":   summary.Frames,
			"media_height":   summary.Height,
			"media_width":    summary.Width,
			"author_data":    authorJson,
			"review_data":    reviewJson,
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
