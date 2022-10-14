package controller

import (
	"fmt"
	"gitee.com/uni-minds/bridge_pacs/tools"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// region define

type screenListCallback struct {
	Data       []screenSeriesDetail `json:"data"`
	ItemsCount int                  `json:"itemsCount"`
}

type screenSeriesDetail struct {
	PatientId       string                 `json:"patient_id"`
	SeriesId        string                 `json:"series_id,omitempty"`
	StudiesId       string                 `json:"studies_id,omitempty"`
	Memo            string                 `json:"memo,omitempty"`
	StudiesMemo     string                 `json:"studies_memo,omitempty"`
	Author          string                 `json:"author"`
	Review          string                 `json:"reviewer"`
	Progress        string                 `json:"progress"`
	StudyDatetime   string                 `json:"studies_datetime"`
	RecordDatetime  string                 `json:"record_datetime"`
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

type ValueString struct {
	Value string
}

type ValueBool struct {
	Value bool
}

type ValueInt struct {
	Value int
}

// endregion

func ScreenGet(ctx *gin.Context) {
	_, exists := ctx.Get("uid")
	if !exists {
		return
	}

	ps := database.BridgeGetPacsServerHandler()
	switch ctx.Query("action") {
	case "getlist":
		var page, row, count int

		gid, err := strconv.Atoi(ctx.Query("gid"))
		if err != nil {
			log("e", err.Error())
			ctx.JSON(http.StatusOK, FailReturn(1000, "wrong group index"))
			return
		}

		index := 0

		if page, _ = strconv.Atoi(ctx.Query("page")); page < 1 {
			page = 1
		}

		if count, _ = strconv.Atoi(ctx.Query("row")); row <= 0 {
			row = 20
		}

		if count, err = strconv.Atoi(ctx.Query("count")); err != nil {
			count = -1
		}

		index = (page - 1) * row

		//studiesInfoList, err := ps.FindStudiesByGroupId(gid, index, count)
		//if err != nil {
		//	ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		//	return
		//}

		studiesIds := module.GroupGetDicom(gid)

		countTotal := len(studiesIds)
		//countTotal := int(ps.CountStudiesByGroupId(gid))
		countRemain := countTotal - index - count
		if countRemain < 0 {
			countRemain = 0
		}

		src := ctx.Query("src")
		switch src {
		case "ui", "UI":
			seriesDetails := make([]screenSeriesDetail, 0)

			//for _, studiesInfo := range studiesInfoList {

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
					seriesDetail.Progress = module.ProgressQueryString(studiesInfo.LabelProgress)
					seriesDetail.StudyDatetime = tools.TimeDecode(studiesInfo.StudyDatetime).Format("2006-01-02 15:04")
					seriesDetail.RecordDatetime = tools.TimeDecode(studiesInfo.RecordDatetime).Format("2006-01-02 15:04")
					seriesDetails = append(seriesDetails, seriesDetail)
				}

			}
			callback := screenListCallback{
				Data:       seriesDetails,
				ItemsCount: countRemain,
			}

			ctx.JSON(http.StatusOK, SuccessReturn(callback))

			//default:
			//	_ = database.UserSetLastStatus(uid, gid, page)
			//	ctx.JSON(http.StatusOK, SuccessReturn(studiesInfoList))
		}

	case "getlock":
		seriesId := ctx.Query("series_id")
		_, uid := CookieValidUid(ctx)

		studiesId, err := module.PacsGetStudiesIdFromSeriesId(seriesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			return
		}

		//fmt.Println("studiesId",studiesId)
		authorUid, err := module.PacsGetStudiesInfoAuthor(studiesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			return
		}

		reviewUid, err := module.PacsGetStudiesInfoReview(studiesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			return
		}

		fmt.Printf("uid org: author -> %d review -> %d\nuid now: %d\n", authorUid, reviewUid, uid)

		if testIsMaster(uid) {
			ctx.JSON(http.StatusOK, SuccessReturn(uid))
		} else if authorUid == 0 || authorUid == uid || reviewUid == uid {
			ctx.JSON(http.StatusOK, SuccessReturn(uid))
		} else {
			ctx.JSON(http.StatusOK, FailReturn(300, uid))
		}
	}
}

func testIsMaster(uid int) bool {
	switch uid {
	case 0, 1: // admin
		return true
	case 4, 6, 26: // master
		return true
	default:
		return false
	}
}

func ScreenPost(ctx *gin.Context) {
	action := ctx.Query("action")
	_, uid := CookieValidUid(ctx)
	switch action {
	case "sync":
		//module.PacsSync("192.168.3.101:8080")
		dh := database.BridgeGetPacsDatabaseHandler()
		dh.Sync(true, true, 50)
		ctx.JSON(http.StatusOK, SuccessReturn("Sync finish"))

	case "author":
		//seriesId := ctx.Query("series_id")
		studiesId := ctx.Query("studies_id")

		pi := database.BridgeGetPacsServerHandler()
		//seriesInfo, err := pi.GetSeries(seriesId)
		//if err != nil {
		//	ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		//	return
		//}

		//studiesId := seriesInfo.StudiesId
		studiesInfo, err := pi.FindStudiesById(studiesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
			return
		}

		if studiesInfo.LabelUidAuthor != 0 && studiesInfo.LabelUidAuthor != uid {
			log("e", "no permission:", uid)
			ctx.JSON(http.StatusOK, FailReturn(404, "no permission"))
			return
		}

		switch ctx.Query("selector") {

		case "submit":
			err = pi.StudiesUpdateLabelProgress(studiesId, 2)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(1))
			}
			return

		case "full":

		default:

		}
	}
}

func ScreenGetStudiesOperation(ctx *gin.Context) {
	studiesId := ctx.Param("studiesId")
	op := ctx.Param("operation")

	si := database.BridgeGetPacsServerHandler()
	si.ShowHidden = true
	si.ShowDelete = true
	info, err := si.FindStudiesById(studiesId)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		return
	}

	var value int

	switch op {
	case "hidden", "hide":
		value = info.DbHidden

	case "delete":
		value = info.DbDelete

	default:
		ctx.JSON(http.StatusOK, FailReturn(404, "operation unknown"))
		return
	}

	switch value {
	case 1:
		ctx.JSON(http.StatusOK, SuccessReturn(true))

	default:
		ctx.JSON(http.StatusOK, SuccessReturn(false))
	}
}
func ScreenPostStudiesOperation(ctx *gin.Context) {
	studiesId := ctx.Param("studiesId")
	op := ctx.Param("operation")

	si := database.BridgeGetPacsServerHandler()

	var data ValueBool

	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		return
	}

	switch op {
	case "hidden", "hide":
		err = si.StudiesSetTagHidden(studiesId, data.Value)

	case "delete":
		err = si.StudiesSetTagDelete(studiesId, data.Value)

	default:
		fmt.Println(studiesId, op, data)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(data))
	}
}

func ScreenGetSeriesOperation(ctx *gin.Context) {
	studiesId := ctx.Param("studiesId")
	seriesId := ctx.Param("seriesId")
	op := ctx.Param("operation")
	fmt.Println(op, studiesId, seriesId)

	switch op {
	case "memo":
		memo, err := module.PacsGetSeriesMemo(seriesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(memo))
		}

	case "info":
		info, err := module.PacsGetSeriesInfo(seriesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(info))
		}

	case "details":
		details, err := ScreenConvertDatabaseToScreenSeriesId(seriesId, true)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(details))
		}

	default:
		ctx.JSON(http.StatusOK, FailReturn(404, "unknown operation"))
	}
}
func ScreenPostSeriesOperation(ctx *gin.Context) {
	studiesId := ctx.Param("studiesId")
	seriesId := ctx.Param("seriesId")
	_, uid := CookieValidUid(ctx)
	op := ctx.Param("operation")
	fmt.Println(uid, op, studiesId, seriesId)

	switch op {
	case "submit":
		action := ctx.Query("action")

		pi := database.BridgeGetPacsServerHandler()
		switch action {
		case "author":
			info, _ := pi.FindStudiesById(studiesId)
			nextProgress := 2
			switch info.LabelProgress {
			case 4:
				nextProgress = 5
			}
			err := pi.StudiesUpdateLabelProgress(studiesId, nextProgress) //提交待审核
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(nextProgress))
			}

		case "review_reject":
			nextProgress := 4
			if err := pi.StudiesUpdateLabelUidReview(studiesId, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			}
			if err := pi.StudiesUpdateLabelProgress(studiesId, nextProgress); err != nil { // 审核拒绝
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(module.ProgressQueryString(nextProgress)))
			}

			//pi.UpdateStudiesLabelUidReview()

		case "review_approve":
			nextProgress := 7 // 审核拒绝
			if err := pi.StudiesUpdateLabelUidReview(studiesId, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			}
			if err := pi.StudiesUpdateLabelProgress(studiesId, nextProgress); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(module.ProgressQueryString(nextProgress)))
			}

		default:
			ctx.JSON(http.StatusOK, FailReturn(404, "unknown operation"))

		}

	case "memo":
		var data ValueString
		err := ctx.BindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else if err = module.PacsSetSeriesMemo(seriesId, data.Value); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(data.Value))
		}

	default:
		ctx.JSON(http.StatusOK, FailReturn(404, "unknown operation"))
	}
}

func ScreenGetInstanceOperation(ctx *gin.Context) {
	studiesId := ctx.Param("studiesId")
	seriesId := ctx.Param("seriesId")
	instanceId := ctx.Param("instanceId")
	_, uid := CookieValidUid(ctx)
	op := ctx.Param("operation")
	pi := database.BridgeGetPacsServerHandler()
	fmt.Println(uid, op, studiesId, seriesId, instanceId)

	switch op {
	case "info":
		info, err := pi.FindInstanceByIdLocal(instanceId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(info))
		}
	}
}
func ScreenPostInstanceOperation(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	studiesId := ctx.Param("studiesId")
	instanceId := ctx.Param("instanceId")
	op := ctx.Param("operation")

	pi := database.BridgeGetPacsServerHandler()

	switch op {
	case "submit":
		var data screenInstanceScreen
		var err error
		if err = ctx.BindJSON(&data); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			return
		}

		switch ctx.Query("action") {
		default:
			if err = pi.StudiesUpdateLabelUidAuthor(studiesId, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				return
			} else if err = pi.StudiesUpdateLabelProgress(studiesId, 1); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				return
			}

		case "review":
			if err = pi.StudiesUpdateLabelUidReview(studiesId, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				return
			} else if err = pi.StudiesUpdateLabelProgress(studiesId, 5); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				return
			}
		}

		switch data.Selector {
		case "view":
			err = pi.InstanceUpdateLabelTag(instanceId, "label_view", data.Value)
		case "diagnose":
			err = pi.InstanceUpdateLabelTag(instanceId, "label_diagnose", data.Value)
		case "interfere":
			err = pi.InstanceUpdateLabelTag(instanceId, "label_interfere", data.Value)
		}

		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(data.Value))
		}
	}
}

func ScreenConvertDatabaseToScreenSeriesId(seriesId string, includeInstanceDetails bool) (detail screenSeriesDetail, err error) {
	ps := database.BridgeGetPacsServerHandler()

	seriesInfo, err := ps.FindSeriesByIdLocal(seriesId)
	if err != nil {
		log("E", err.Error())
		return detail, err
	}

	details := make([]screenInstanceDetail, 0)

	instanceIds := strings.Split(seriesInfo.IncludeInstances, "|")
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
		SeriesId:        seriesInfo.SeriesId,
		Memo:            seriesInfo.LabelMemo,
		StudiesId:       seriesInfo.StudiesId,
		StudiesMemo:     "",
		InstanceDetails: details,
		InstanceCount:   len(instanceIds),
	}

	return detail, nil
}

func ScreenConvertDatabaseToInstanceId(instanceId string) (detail screenInstanceDetail, err error) {
	ps := database.BridgeGetPacsServerHandler()

	instanceInfo, err := ps.FindInstanceByIdLocal(instanceId)
	if err != nil {
		log("E", err.Error())
		return detail, err
	}

	detail = screenInstanceDetail{
		InstanceId:     instanceInfo.InstanceId,
		Memo:           instanceInfo.LabelMemo,
		Frames:         instanceInfo.Frames,
		LabelView:      instanceInfo.LabelView,
		LabelDiagnose:  instanceInfo.LabelDiagnose,
		LabelInterfere: instanceInfo.LabelInterfere,
	}

	return detail, nil
}

func ScreenDelete(ctx *gin.Context) {
	studiesId := ctx.Query("studies_id")

	_, uid := CookieValidUid(ctx)

	pi := database.BridgeGetPacsServerHandler()
	info, err := pi.FindStudiesById(studiesId)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, "异常"))
		return
	}

	if info.LabelUidAuthor != uid {
		ctx.JSON(http.StatusOK, FailReturn(400, "非标注人，禁止操作"))
		return
	}

	pi.StudiesUpdateLabelProgress(studiesId, 0)
	pi.StudiesUpdateLabelUidAuthor(studiesId, 0)

	ctx.JSON(http.StatusOK, SuccessReturn("OK"))
}
