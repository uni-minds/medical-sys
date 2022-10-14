/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: apiLabelsysGetStream.go
 * Date: 2022/09/15 13:05:15
 */

// ECode = 40200x

package controller

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func LabelsysPostStream(ctx *gin.Context) {
	// apiGroup.GET("/labelsys/stream/index/:index/:class/:op", controller.LabelsysGetStream)
	uid := -1
	if uidi, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = uidi.(int)
	}

	mediaUUID := ctx.Param("index")
	opClass := ctx.Param("class")
	opType := ctx.Param("op")

	info, err := module.MediaGet(mediaUUID)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(402009, err.Error()))
		return
	}
	switch opClass {
	case "label":
		var postLabelData LabelData
		if err = ctx.BindJSON(&postLabelData); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(402012, err.Error()))
			return
		}

		switch opType {
		case "memo":
			if info.Memo != postLabelData.Data {
				log("d", "change memo", info.Memo, "->", postLabelData.Data)
				if err = module.MediaSetMemo(mediaUUID, postLabelData.Data); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(402013, err.Error()))
					return
				}
			}

			ctx.JSON(http.StatusOK, SuccessReturn(postLabelData.Data))
			return

		case "author":
			switch ctx.Query("do") {
			case "submit":
				// 提交
				if err = module.ProcessAuthorSubmit(mediaUUID, uid, postLabelData.Data); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				}
				return

			default:
				if err = module.ProcessAuthorData(mediaUUID, uid, postLabelData.Data); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(402014, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn(fmt.Sprintf("同步成功 @ %s", time.Now().Format(global.TimeFormat))))
				}
				return
			}

		case "review":
			switch ctx.Query("do") {
			case "confirm":
				if err = module.ProcessReviewConfirm(mediaUUID, postLabelData.Data, uid); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				}

			case "reject":
				if err = module.ProcessReviewReject(mediaUUID, postLabelData.Data, uid); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				}

			case "revoke":
				if err = module.ProcessReviewRevoke(mediaUUID, postLabelData.Data, uid, postLabelData.Admin); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				}

			default:
				if err = module.LabelUpdateReview(postLabelData.Data, mediaUUID, uid); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(402014, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn("ok"))
				}
				return
			}

		default:
			ctx.JSON(http.StatusOK, FailReturn(402011, "unknown label operation"))
			return
		}

	default:
		ctx.JSON(http.StatusOK, FailReturn(402010, "unknown operate class"))
	}
}
