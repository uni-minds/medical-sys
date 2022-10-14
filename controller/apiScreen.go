package controller

import (
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
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
type screenSubmitData struct {
	Operate     string `json:"operate,omitempty"`
	StudiesId   string `json:"studies_id,omitempty"`
	SeriesId    string `json:"series_id,omitempty"`
	InstanceId  string `json:"instance_id,omitempty"`
	SubmitLevel string `json:"submit_level,omitempty"`
	Info        struct {
		Selector string `json:"selector,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"info"`
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
type ValueAdmin struct {
	Admin string
}

// endregion

func PostStudiesOperation(ctx *gin.Context) {
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

func SeriesGetOperation(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	studiesId := ctx.Param("studiesId")
	seriesId := ctx.Param("seriesId")
	op := ctx.Param("operation")
	log("d", op, studiesId, seriesId)

	switch op {
	case "screen_getlock":
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

		log("d", fmt.Sprintf("uid org: author -> %d review -> %d\nuid now: %d\n", authorUid, reviewUid, uid))

		if testIsMaster(uid) {
			ctx.JSON(http.StatusOK, SuccessReturn(uid))
		} else if authorUid == 0 || authorUid == uid || reviewUid == uid {
			ctx.JSON(http.StatusOK, SuccessReturn(uid))
		} else {
			ctx.JSON(http.StatusOK, FailReturn(300, uid))
		}

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
func SeriesPostOperation(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	studiesId := ctx.Param("studiesId")
	seriesId := ctx.Param("seriesId")
	op := ctx.Param("operation")
	log("d", uid, op, studiesId, seriesId)

	switch op {
	case "screen_submit":
		var data screenSubmitData
		if err := ctx.BindJSON(&data); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(403, "data error"))
		} else if status, err := ParseScreenData(uid, studiesId, seriesId, data); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(403, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(status))
		}

	case "memo":
		var data LabelData
		err := ctx.BindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else if err = module.PacsSetSeriesMemo(seriesId, data.Data); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(data.Data))
		}

	default:
		ctx.JSON(http.StatusOK, FailReturn(404, "unknown operation"))
	}
}

func SeriesDelOperation(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	studiesId := ctx.Param("studiesId")
	seriesId := ctx.Param("seriesId")
	operation := ctx.Param("operation")
	log("d", uid, studiesId, seriesId, operation, "delete")

	var value ValueAdmin
	if err := ctx.BindJSON(&value); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return
	}

	pi := database.BridgeGetPacsServerHandler()

	switch operation {
	// 删除挑图数据
	case "screen":
		info, err := pi.FindStudiesById(studiesId)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, "异常"))
			return
		}

		if value.Admin != "" && value.Admin != global.DefAdminPassword && info.LabelUidAuthor != uid {
			ctx.JSON(http.StatusOK, FailReturn(400, "非标注人，禁止操作"))
			return
		}

		pi.StudiesUpdateLabelProgress(studiesId, 0)
		pi.StudiesUpdateLabelUidAuthor(studiesId, 0)

		ctx.JSON(http.StatusOK, SuccessReturn("OK"))

	// 删除原始数据
	case "raw":
		err := pi.RemoveSeries(studiesId, seriesId, true)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn("Deleted"))
		}
	}
}

func InstanceGetOperation(ctx *gin.Context) {

	studiesId := ctx.Param("studiesId")
	seriesId := ctx.Param("seriesId")
	instanceId := ctx.Param("instanceId")
	op := ctx.Param("operation")
	pi := database.BridgeGetPacsServerHandler()

	log("d", "user", op, studiesId, seriesId, instanceId)

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
func ParseScreenData(uid int, studies_id, series_id string, data screenSubmitData) (status string, err error) {
	pi := database.BridgeGetPacsServerHandler()

	switch data.Operate {
	case "instance_set":
		switch data.SubmitLevel {
		case "author":
			if err = pi.StudiesUpdateLabelUidAuthor(studies_id, uid); err != nil {
				return "", err
			} else if err = pi.StudiesUpdateLabelProgress(studies_id, 1); err != nil {
				return "", err
			}

		case "review":
			if err = pi.StudiesUpdateLabelUidReview(studies_id, uid); err != nil {
				return "", err
			} else if err = pi.StudiesUpdateLabelProgress(studies_id, 5); err != nil {
				return "", err
			}
		}

		switch data.Info.Selector {
		case "view":
			err = pi.InstanceUpdateLabelTag(data.InstanceId, "label_view", data.Info.Value)
		case "diagnose":
			err = pi.InstanceUpdateLabelTag(data.InstanceId, "label_diagnose", data.Info.Value)
		case "interfere":
			err = pi.InstanceUpdateLabelTag(data.InstanceId, "label_interfere", data.Info.Value)
		}

		return data.Info.Value, nil

	case "author_series_submit":
		info, _ := pi.FindStudiesById(studies_id)
		nextProgress := 2
		switch info.LabelProgress {
		case 4:
			nextProgress = 5
		}
		//提交待审核
		if err = pi.StudiesUpdateLabelProgress(studies_id, nextProgress); err != nil {
			return "", err
		} else {
			return module.ProgressQuery(nextProgress), nil
		}

	case "review_series_reject":
		nextProgress := 4
		if err = pi.StudiesUpdateLabelUidReview(studies_id, uid); err != nil {
			return "", err
		} else if err = pi.StudiesUpdateLabelProgress(studies_id, nextProgress); err != nil { // 审核拒绝
			return "", err
		} else {
			return module.ProgressQuery(nextProgress), nil
		}

	case "review_series_approve":
		nextProgress := 7 // 审核拒绝
		if err = pi.StudiesUpdateLabelUidReview(studies_id, uid); err != nil {
			return "", err
		} else if err = pi.StudiesUpdateLabelProgress(studies_id, nextProgress); err != nil {
			return "", err
		} else {
			return module.ProgressQuery(nextProgress), nil
		}

	default:
		return "", errors.New(fmt.Sprint("unknown submit operate:", data.Operate))
	}
}
