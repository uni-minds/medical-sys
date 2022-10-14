package database

import (
	"errors"
	"fmt"
	"log"
	"uni-minds.com/liuxy/medical-sys/logger"
	"uni-minds.com/liuxy/medical-sys/tools"
)

const TablePacsStudies = "pacs_studies"
const TablePacsSeries = "pacs_series"
const TablePacsInstances = "pacs_instances"

type PacsInstanceDetail struct {
	Server       string
	AET          string
	StudiesID    string
	SeriesID     string
	InstanceID   string
	Frames       int
	Width        int
	Height       int
	InstanceType string
}

type InfoStudies struct {
	StudiesID       string `gorose:"studies_id"`
	IncludeSeriesID string `gorose:"include_series_id"`
	LabelInfo       string `gorose:"label_info"`
}

var SqlInitDbStudies = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"studies_id" TEXT NOT NULL PRIMARY KEY DEFAULT "",
	"include_series_id" TEXT NOT NULL DEFAULT "",
	"label_info" TEXT NOT NULL DEFAULT "")`, TablePacsStudies)

func (*InfoStudies) TableName() string {
	return TablePacsStudies
}

type InfoSeries struct {
	SeriesID          string `gorose:"series_id"`
	IncludeInstanceID string `gorose:"include_instance_id"`
	LabelAuthorUid    int    `gorose:"label_author_uid"`
	LabelReviewUid    int    `gorose:"label_review_uid"`
	LabelProgress     int    `gorose:"label_progress"`
	LabelInfo         string `gorose:"label_info"`
}

var SqlInitDbSeries = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"series_id" TEXT NOT NULL DEFAULT "",
	"include_instance_id" TEXT NOT NULL DEFAULT "",
	"label_author_uid" INTEGER DEFAULT 0, 
	"label_review_uid" INTEGER DEFAULT 0,
	"label_progress" INTEGER DEFAULT 0,
	"label_info" TEXT NOT NULL DEFAULT "")`, TablePacsSeries)

func (*InfoSeries) TableName() string {
	return TablePacsSeries
}

type InfoInstance struct {
	InstanceID     string `gorose:"instance_id"`
	Server         string `gorose:"server"`
	AET            string `gorose:"aet"`
	StudiesID      string `gorose:"studies_id"`
	SeriesID       string `gorose:"series_id"`
	DisplayName    string `gorose:"display_name"`
	PathCache      string `gorose:"path_cache"`
	Frames         int    `gorose:"frames"`
	Width          int    `gorose:"width"`
	Height         int    `gorose:"height"`
	InstanceType   string `gorose:"instance_type"`
	LabelInfo      string `gorose:"label_info"`
	Memo           string `gorose:"memo"`
	LabelView      string `gorose:"label_view"`
	LabelDiagnose  string `gorose:"label_diagnose"`
	LabelInterfere string `gorose:"label_interfere"`
}

var SqlInitDbInstance = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"instance_id" TEXT PRIMARY KEY,
	"server" TEXT NOT NULL DEFAULT "",
	"aet" TEXT NOT NULL default "DCM4CHEE",
	"studies_id" TEXT NOT NULL DEFAULT "",
	"series_id" TEXT NOT NULL DEFAULT "",
	"display_name" TEXT NOT NULL DEFAULT "",
	"path_cache" TEXT NOT NULL DEFAULT "",
	"frames" INTEGER NOT NULL default 0,
	"width" INTEGER NOT NULL default 0,
	"height" INTEGER NOT NULL default 0,
	"instance_type" INTERGET NOT NULL default 0,
	"label_info" TEXT NOT NULL DEFAULT "",
	"label_view" INTEGER NOT NULL default 0,
	"label_diagnose" INTEGER NOT NULL DEFAULT 0,
	"label_interfere" INTEGER NOT NULL DEFAULT 0,
	"memo" TEXT NOT NULL default "")`, TablePacsInstances)

func (*InfoInstance) TableName() string {
	return TablePacsInstances
}

func initPacsDB() {
	if _, err := DB().Execute(SqlInitDbStudies); err != nil {
		log.Panic(err.Error())
	}

	if _, err := DB().Execute(SqlInitDbSeries); err != nil {
		log.Panic(err.Error())
	}

	if _, err := DB().Execute(SqlInitDbInstance); err != nil {
		log.Panic(err.Error())
	}
}

func PacsStudiesGetInfo(studiesId string) (info InfoStudies, err error) {
	err = DB().Table(&info).Where("studies_id", "=", studiesId).Select()
	if info.StudiesID != studiesId {
		return info, errors.New("studies id not found")
	} else {
		return info, nil
	}
}
func PacsStudiesCreate(studiesId, seriesId string) (err error) {
	str, _ := tools.StringCompress([]string{seriesId})

	data := InfoStudies{
		StudiesID:       studiesId,
		IncludeSeriesID: str,
	}
	//fmt.Println("create pacs data",data)
	_, err = DB().Table(TablePacsStudies).Data(data).Insert()
	return err
}
func PacsStudiesInsert(studiesId, seriesId string) (err error) {
	var studiesIDs []string
	studiesInfo, _ := PacsStudiesGetInfo(studiesId)
	if studiesIDs, err = tools.StringDecompress(studiesInfo.IncludeSeriesID); err != nil {
		fmt.Println("DB_E1:", err.Error(), studiesInfo.IncludeSeriesID)
		return err
	}

	for _, sid := range studiesIDs {
		if sid == seriesId {
			return nil
		}
	}

	studiesIDs = append(studiesIDs, seriesId)
	if str, err := tools.StringCompress(studiesIDs); err != nil {
		return err
	} else {
		studiesInfo.IncludeSeriesID = str
		_, err = DB().Table(TablePacsSeries).Data(studiesInfo).Where("studies_id", "=", studiesId).Update()
		return err
	}
}

func PacsSeriesGetInfo(seriesId string) (info InfoSeries, err error) {
	err = DB().Table(&info).Where("series_id", "=", seriesId).Select()
	if info.SeriesID != seriesId {
		return info, errors.New("series id not found")
	} else {
		return info, nil
	}
}
func PacsSeriesCreate(seriesId, instanceId string) (err error) {
	str, _ := tools.StringCompress([]string{instanceId})
	data := InfoSeries{
		SeriesID:          seriesId,
		IncludeInstanceID: str,
		LabelAuthorUid:    0,
		LabelReviewUid:    0,
		LabelProgress:     0,
		LabelInfo:         "",
	}
	_, err = DB().Table(TablePacsSeries).Data(data).Insert()
	return err
}
func PacsSeriesInsert(seriesId, instanceId string) (err error) {
	var instanceIds []string
	info, _ := PacsSeriesGetInfo(seriesId)
	if instanceIds, err = tools.StringDecompress(info.IncludeInstanceID); err != nil {
		return err
	}

	for _, sid := range instanceIds {
		if sid == instanceId {
			return nil
		}
	}

	instanceIds = append(instanceIds, instanceId)
	if str, err := tools.StringCompress(instanceIds); err != nil {
		return err
	} else {
		info.IncludeInstanceID = str
		_, err = DB().Table(TablePacsSeries).Data(info).Where("series_id", "=", seriesId).Update()
		return err
	}
}
func PacsSeriesUpdate(data InfoSeries) (err error) {
	if data.SeriesID == "" {
		return errors.New("empty instance id")
	}

	if _, err = DB().Table(TablePacsSeries).Data(data).Where("series_id", "=", data.SeriesID).Update(); err != nil {
		logger.Write("DB", "E", err.Error())
		return err
	}
	return nil
}
func PacsSeriesUpdateAuthor(seriesId string, uid, progress int) (err error) {
	info, err := PacsSeriesGetInfo(seriesId)
	if err != nil {
		return err
	}

	info.LabelAuthorUid = uid
	info.LabelProgress = progress

	return PacsSeriesUpdate(info)
}
func PacsSeriesUpdateReview(seriesId string, uid, progress int) (err error) {
	info, err := PacsSeriesGetInfo(seriesId)
	if err != nil {
		return err
	}

	info.LabelProgress = progress
	info.LabelReviewUid = uid

	return PacsSeriesUpdate(info)
}

func PacsInstanceGetInfo(instanceId string) (info InfoInstance, err error) {
	err = DB().Table(&info).Where("instance_id", "=", instanceId).Select()
	if info.InstanceID != instanceId {
		return info, errors.New("instance id not found")
	} else {
		return info, nil
	}
}
func PacsInstanceCreate(data PacsInstanceDetail) (err error) {
	if _, err = PacsInstanceGetInfo(data.InstanceID); err != nil {
		instanceData := InfoInstance{
			InstanceID:   data.InstanceID,
			Server:       data.Server,
			AET:          data.AET,
			StudiesID:    data.StudiesID,
			SeriesID:     data.SeriesID,
			DisplayName:  data.InstanceID,
			PathCache:    "",
			Frames:       data.Frames,
			Width:        data.Width,
			Height:       data.Height,
			InstanceType: data.InstanceType,
			LabelInfo:    "",
			Memo:         "",
		}

		if _, err = DB().Table(TablePacsInstances).Data(instanceData).Insert(); err != nil {
			logger.Write("DB", "E", err.Error())
		}
	} else {
		logger.Write("DB", "w", fmt.Sprintf("Instance already exist: %s", data.InstanceID))
	}

	if _, err = PacsSeriesGetInfo(data.SeriesID); err != nil {
		if err = PacsSeriesCreate(data.SeriesID, data.InstanceID); err != nil {
			logger.Write("DB", "E", err.Error())
		}
	} else if err = PacsSeriesInsert(data.SeriesID, data.InstanceID); err != nil {
		logger.Write("DB", "E", err.Error())

	}

	if _, err = PacsStudiesGetInfo(data.StudiesID); err != nil {
		err = PacsStudiesCreate(data.StudiesID, data.SeriesID)
	} else {
		err = PacsStudiesInsert(data.StudiesID, data.SeriesID)
	}

	if err != nil {
		fmt.Println("E11", err.Error())
	}
	return nil
}
func PacsInstanceUpdate(data InfoInstance) (err error) {
	if data.InstanceID == "" {
		return errors.New("empty instance id")
	}

	if _, err = DB().Table(TablePacsInstances).Data(data).Where("instance_id", "=", data.InstanceID).Update(); err != nil {
		logger.Write("DB", "E", err.Error())
		return err
	}
	return nil
}
func PacsInstanceUpdateLabel(instanceId, view, diagnose, interfere string) (err error) {
	info, err := PacsInstanceGetInfo(instanceId)
	if err != nil {
		return err
	}

	if view != "" {
		info.LabelView = view
	}

	if diagnose != "" {
		info.LabelDiagnose = diagnose
	}

	if interfere != "" {
		info.LabelInterfere = interfere
	}

	return PacsInstanceUpdate(info)
}

func PacsGetAllStudiesIds() (studiesIds []string, err error) {
	var infos []InfoStudies
	if err = DB().Table(&infos).Select(); err != nil {
		return nil, err
	} else {
		for _, info := range infos {
			studiesIds = append(studiesIds, info.StudiesID)
		}
		return studiesIds, nil
	}
}
