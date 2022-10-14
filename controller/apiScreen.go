package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/module"
	"uni-minds.com/liuxy/medical-sys/tools"
)

type screenListCallback struct {
	Data       []screenSeriesDetail `json:"data"`
	ItemsCount int                  `json:"itemsCount"`
}

type screenSeriesDetail struct {
	SeriesId        string                 `json:"series_id,omitempty"`
	Memo            string                 `json:"memo,omitempty"`
	StudiesId       string                 `json:"studies_id,omitempty"`
	StudiesMemo     string                 `json:"studies_memo,omitempty"`
	Author          string                 `json:"author"`
	Review          string                 `json:"reviewer"`
	Progress        int                    `json:"progress"`
	InstanceDetails []screenInstanceDetail `json:"instance_details,omitempty"`
	InstanceCount   int                    `json:"instance_count"`
}

type screenInstanceDetail struct {
	InstanceId     string `json:"instance_id,omitempty"`
	Memo           string `json:"memo,omitempty"`
	Frames         int    `json:"frames"`
	LabelView      string `json:"label_view,omitempty"`      //切面
	LabelDiagnose  string `json:"label_diagnose,omitempty"`  //是否正常
	LabelInterfere string `json:"label_interfere,omitempty"` //存在测量干扰
}

type screenInstanceScreen struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

func ScreenGet(ctx *gin.Context) {
	_, uid := CookieValidUid(ctx)
	action := ctx.Query("action")
	switch action {
	case "getlist":
		var page, count int
		var studiesIds []string

		gid, err := strconv.Atoi(ctx.Query("gid"))
		if err != nil {
			log("e", err.Error())
			ctx.JSON(http.StatusOK, FailReturn(1000, "wrong group index"))
			return
		}

		if page, _ = strconv.Atoi(ctx.Query("page")); page <= 0 {
			page = 1
		}

		if count, _ = strconv.Atoi(ctx.Query("count")); count <= 0 {
			count = 20
		}

		index := (page - 1) * count

		studiesIds, err = database.GroupGetPacsStudiesIds(gid)
		if err != nil {
			log("E", err.Error())
			ctx.JSON(http.StatusOK, FailReturn(1002, err.Error()))
		}

		if index >= len(studiesIds) {
			ctx.JSON(http.StatusOK, FailReturn(1002, "out of range"))
			return
		}

		studiesIds = studiesIds[index:]

		totalRemain := len(studiesIds)

		if count > totalRemain {
			count = totalRemain
		}

		src := ctx.Query("src")
		switch src {
		case "ui", "UI":
			seriesDetails := make([]screenSeriesDetail, 0)

			for _, studiesId := range studiesIds[0:count] {
				studiesInfo, err := database.PacsStudiesGetInfo(studiesId)
				if err != nil {
					log("E", err.Error())
					continue
				}

				seriesIds, err := tools.StringDecompress(studiesInfo.IncludeSeriesID)
				if err != nil {
					log("E", err.Error())
					continue
				}

				for _, seriesId := range seriesIds {
					seriesDetail, err := ScreenConvertDatabaseToScreenSeriesId(seriesId, false)
					if err != nil {
						log("e", "series", err.Error())
						continue
					}
					seriesDetail.StudiesId = studiesInfo.StudiesID
					seriesDetail.StudiesMemo = studiesInfo.LabelInfo
					seriesDetails = append(seriesDetails, seriesDetail)
				}

			}
			callback := screenListCallback{
				Data:       seriesDetails,
				ItemsCount: totalRemain,
			}

			ctx.JSON(http.StatusOK, SuccessReturn(callback))

		default:
			if studiesIds, err = database.GroupGetPacsStudiesIds(gid); err != nil {
				log("E", err.Error())
				ctx.JSON(http.StatusOK, FailReturn(1001, err.Error()))
				return
			}

			if len(studiesIds) <= index {
				ctx.JSON(http.StatusOK, FailReturn(1002, "list index"))
				return
			}

			_ = database.UserSetLastStatus(uid, gid, page)
			studiesIds = studiesIds[index:]

			ctx.JSON(http.StatusOK, studiesIds)
		}

	case "getdetails":
		seriesId := ctx.Query("series_id")

		details, err := ScreenConvertDatabaseToScreenSeriesId(seriesId, true)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(details))
		}

	case "getinstance":
		instanceId := ctx.Query("instance_id")
		//thumbs := ctx.Query("thumbs")
		info, err := database.PacsInstanceGetInfo(instanceId)
		if err != nil {
			ctx.JSON(http.StatusOK, info)
		}

	case "getmedia":
		instanceId := ctx.Query("instance_id")
		regen := ctx.Query("force_regen")
		thumbSize := ctx.Query("thumb")
		video := ctx.Query("video")
		var bs []byte
		var err error
		switch thumbSize {
		case "300":
			bs, err = module.PacsGetInstanceThumb(instanceId)

		default:
			if video == "true" {
				bs, _, err = module.PacsGetInstanceImage(instanceId, regen != "")
			} else {
				bs, _, err = module.PacsGetInstanceImage(instanceId, regen != "")
			}
		}
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.Writer.Write(bs)
		}

	case "getlock":
		seriesId := ctx.Query("series_id")
		_, uid := CookieValidUid(ctx)

		info, err := database.PacsSeriesGetInfo(seriesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			return
		}

		if info.LabelAuthorUid == 0 || uid == info.LabelAuthorUid || uid == info.LabelReviewUid {
			ctx.JSON(http.StatusOK, SuccessReturn(uid))
		} else {
			ctx.JSON(http.StatusOK, FailReturn(300, uid))
		}
	}
}

func ScreenConvertDatabaseToScreenSeriesId(seriesId string, includeInstanceDetails bool) (detail screenSeriesDetail, err error) {
	seriesInfo, err := database.PacsSeriesGetInfo(seriesId)
	if err != nil {
		log("E", err.Error())
		return detail, err
	}

	details := make([]screenInstanceDetail, 0)

	instanceIds, err := tools.StringDecompress(seriesInfo.IncludeInstanceID)
	if includeInstanceDetails {
		for _, instanceId := range instanceIds {
			instanceDetail, err := ScreenConvertDatabaseToInstanceId(instanceId)
			if err != nil {
				instanceDetail = screenInstanceDetail{
					InstanceId: instanceId,
					Memo:       "无信息",
				}
			}
			details = append(details, instanceDetail)
		}
	}

	detail = screenSeriesDetail{
		SeriesId:        seriesInfo.SeriesID,
		Memo:            seriesInfo.LabelInfo,
		StudiesId:       "",
		StudiesMemo:     "",
		InstanceDetails: details,
		InstanceCount:   len(instanceIds),
	}

	if seriesInfo.LabelAuthorUid > 0 {
		detail.Author = module.UserGetRealname(seriesInfo.LabelAuthorUid)
	}

	if seriesInfo.LabelProgress > 0 {
		detail.Progress = seriesInfo.LabelProgress
	}

	if seriesInfo.LabelReviewUid > 0 {
		detail.Review = module.UserGetRealname(seriesInfo.LabelReviewUid)
	}
	return detail, nil
}

func ScreenConvertDatabaseToInstanceId(instanceId string) (detail screenInstanceDetail, err error) {
	instanceInfo, err := database.PacsInstanceGetInfo(instanceId)
	if err != nil {
		log("E", err.Error())
		return detail, err
	}

	detail = screenInstanceDetail{
		InstanceId:     instanceInfo.InstanceID,
		Memo:           instanceInfo.Memo,
		Frames:         instanceInfo.Frames,
		LabelView:      instanceInfo.LabelView,
		LabelDiagnose:  instanceInfo.LabelDiagnose,
		LabelInterfere: instanceInfo.LabelInterfere,
	}

	return detail, nil
}

func ScreenPost(ctx *gin.Context) {
	var err error
	action := ctx.Query("action")
	_, uid := CookieValidUid(ctx)
	switch action {
	case "sync":
		module.PacsSync("192.168.3.101:8080")
		ctx.JSON(http.StatusOK, SuccessReturn("Sync finish"))

	case "author":
		//studies_id := ctx.Query("studies_id")
		series_id := ctx.Query("series_id")
		instanceId := ctx.Query("instance_id")

		switch ctx.Query("selector") {
		case "memo":

		case "submit":
			err = module.PacsSetSeriesAuthorLabel(series_id, uid, 2)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(1))
			}
			return

		case "full":

		default:
			var data screenInstanceScreen

			err = module.PacsSetSeriesAuthorLabel(series_id, uid, 1)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			}

			err = ctx.BindJSON(&data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			}

			switch data.Selector {
			case "view":
				err = module.PacsInstanceUpdateLabel(series_id, instanceId, uid, data.Value, "", "")
			case "diagnose":
				err = module.PacsInstanceUpdateLabel(series_id, instanceId, uid, "", data.Value, "")
			case "interfere":
				err = module.PacsInstanceUpdateLabel(series_id, instanceId, uid, "", "", data.Value)
			}

			if err != nil {
				log("e", err.Error())
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(data.Value))
			}
		}
	}
}

func ScreenDelete(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, SuccessReturn("OK Delete"))
}
