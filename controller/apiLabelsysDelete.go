/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: apiLabelsysGetStream.go
 * Date: 2022/09/15 13:05:15
 */

// ECode = 40200x

package controller

import (
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LabelsysDeleteStream(ctx *gin.Context) {
	// apiGroup.Del("/labelsys/stream/index/:index/:class/:op", controller.LabelsysGetStream)

	mediaUUID := ctx.Param("index")
	opClass := ctx.Param("class")
	opType := ctx.Param("op")

	//info, err := module.MediaGet(mediaUUID)
	//if err != nil {
	//	ctx.JSON(http.StatusOK, FailReturn(402009, err.Error()))
	//	return
	//}
	switch opClass {
	case "label":
		var postLabelData LabelData
		if err := ctx.BindJSON(&postLabelData); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(402012, err.Error()))
			return
		}

		switch opType {
		case "full":
			log("d", "user delete label for media:", mediaUUID)

			if postLabelData.Admin == global.DefAdminPassword {
				if err := module.MediaRemoveLabelAll(mediaUUID); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				}
			} else {
				log("w", "pwd=", global.DefAdminPassword)
				ctx.JSON(http.StatusOK, FailReturn(http.StatusForbidden, "禁止操作"))
			}
			return

		case "review":
			log("d", "user delete review", mediaUUID)
			ctx.JSON(http.StatusOK, SuccessReturn(postLabelData.Data))

		default:
			ctx.JSON(http.StatusOK, FailReturn(402011, "unknown label operation"))
			return
		}

	default:
		ctx.JSON(http.StatusOK, FailReturn(402010, "unknown operate class"))
	}
}
