/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: apiMedialist.go
 * Date: 2022/09/19 08:15:19
 */

package controller

import (
	pacs_tools "gitee.com/uni-minds/bridge-pacs/tools"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// /api/v1/medialist/:groupType/:groupId/:op
func MedialistGet(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	operate := ctx.Param("op")

	gid, err := strconv.Atoi(ctx.Param("groupId"))
	if err != nil {
		return
	}

	switch operate {

	// /api/v1/medialist/:groupType/:groupId/list
	case "list":
		var callback medialistForJqGrid

		page, _ := strconv.Atoi(ctx.Query("page"))
		if page <= 0 {
			page = 1
		}
		rows, _ := strconv.Atoi(ctx.Query("rows"))
		if rows <= 0 {
			rows = 10
		}
		order := ctx.Query("order")
		field := ctx.Query("field")
		view := ctx.Query("view")
		//ignoreProgressCheck := ctx.Query("ignoreProgressCheck")

		index := (page - 1) * rows

		mediaInfos, total, err := module.UserGetGroupContainsWithSelector(uid, gid, view, field, order, index, rows)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(403, err.Error()))
			return
		}

		callback.Records = total
		callback.Page = page
		callback.Total = int(math.Ceil(float64(total) / float64(rows)))

		mdata := make([]mediaInfoForJqGrid, 0)

		for _, mediaSummary := range module.ConvertSummaryFromMediaInfos(mediaInfos) {
			labelUUIDs, err := module.MediaGetLabelUUIDs(mediaSummary.MediaUUID)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(403, "get label uuids:"+err.Error()))
				return
			}

			authorsData := make([]labelInfoForJsGridButton, 0)
			reviewsData := make([]labelInfoForJsGridButton, 0)
			for _, labelUUID := range labelUUIDs {
				if labelSummary, _, _, err := module.LabelGetSummary(labelUUID); err == nil {
					authorsData = append(authorsData, labelInfoForJsGridButton{
						Realname: labelSummary.AuthorRealname,
						Tips:     labelSummary.AuthorTips,
						Status:   labelSummary.AuthorProgress,
						UUID:     labelUUID,
					})

					reviewsData = append(reviewsData, labelInfoForJsGridButton{
						Realname: labelSummary.ReviewRealname,
						Tips:     labelSummary.ReviewTips,
						Status:   labelSummary.ReviewProgress,
						UUID:     labelUUID,
					})
				}
			}

			mdata = append(mdata, mediaInfoForJqGrid{
				Mid:       strconv.Itoa(mediaSummary.Id),
				MediaUUID: mediaSummary.MediaUUID,
				Name:      mediaSummary.DisplayName,
				View:      mediaSummary.Views,
				Duration:  mediaSummary.Duration,
				Frames:    mediaSummary.Frames,
				Authors:   authorsData,
				Reviews:   reviewsData,
				Memo:      mediaSummary.Memo,
			})
		}

		callback.Rows = make([]interface{}, 0)
		for _, row := range mdata {
			callback.Rows = append(callback.Rows, row)
		}
		ctx.JSON(http.StatusOK, callback)

		return

	// /api/v1/medialist/groupType/groupId/screen
	case "screen":
		var callback medialistForJqGrid

		page, _ := strconv.Atoi(ctx.Query("page"))
		if page <= 0 {
			page = 1
		}
		rows, _ := strconv.Atoi(ctx.Query("rows"))
		if rows <= 0 {
			rows = 10
		}
		//order := ctx.Query("order")
		//field := ctx.Query("field")
		//view := ctx.Query("view")
		//ignoreProgressCheck := ctx.Query("ignoreProgressCheck")

		index := (page - 1) * rows

		studiesIds, containType := module.GroupGetContainMedia(gid)
		if containType != "studies_id" {
			ctx.JSON(http.StatusOK, FailReturn(10001, "此组未开通数据筛选功能:"+containType))
			return
		}

		total := len(studiesIds)

		callback.Records = total
		callback.Page = page
		callback.Total = int(math.Ceil(float64(total) / float64(rows)))

		studiesIds = studiesIds[index:]
		if len(studiesIds) > rows {
			studiesIds = studiesIds[:rows]
		}

		ps := database.BridgeGetPacsServerHandler()

		mdata := make([]screenSeriesDetail, 0)

		for _, studiesId := range studiesIds {
			studiesInfo, err := ps.FindStudiesById(studiesId)
			if err != nil {
				log("e", err.Error())
				continue
			}

			seriesIds := strings.Split(studiesInfo.IncludeSeries, "|")

			for _, seriesId := range seriesIds {
				seriesDetail, err := ScreenConvertDatabaseToScreenSeriesId(seriesId, false)
				if err != nil {
					log("e", "series", err.Error())
					continue
				}
				seriesDetail.PatientId = studiesInfo.PatientId
				seriesDetail.StudiesId = studiesInfo.StudiesId
				seriesDetail.StudiesMemo = studiesInfo.LabelMemo
				seriesDetail.Author = module.UserGetRealname(studiesInfo.LabelUidAuthor)
				seriesDetail.Review = module.UserGetRealname(studiesInfo.LabelUidReview)
				seriesDetail.Progress = module.ProgressQuery(studiesInfo.LabelProgress)
				seriesDetail.StudyDatetime = pacs_tools.TimeDecode(studiesInfo.StudyDatetime).Format("2006-01-02 15:04")
				seriesDetail.RecordDatetime = pacs_tools.TimeDecode(studiesInfo.RecordDatetime).Format("2006-01-02 15:04")
				mdata = append(mdata, seriesDetail)
			}
		}

		callback.Rows = make([]interface{}, 0)
		for _, row := range mdata {
			callback.Rows = append(callback.Rows, row)
		}

		ctx.JSON(http.StatusOK, callback)

	// /api/v1/medialist/:groupType/:groupId/view
	case "view":
		ignoreProgressCheck := ctx.Query("ignoreProgressCheck")

		views, err := module.GroupGetContainViews(gid, uid, ignoreProgressCheck != "")
		if err != nil {
			log("e", "group get contain views:", err.Error())
			ctx.JSON(http.StatusOK, FailReturn(403, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(views))
		}
	}
}
