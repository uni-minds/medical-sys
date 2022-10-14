/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: media.go
 */

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"strings"
)

func (*MediaInfo) TableName() string {
	return global.DefaultDatabaseMediaTable
}

type MediaInfo struct {
	Mid            int     `gorose:"mid"`
	DisplayName    string  `gorose:"displayname"`
	Path           string  `gorose:"path"`
	Hash           string  `gorose:"hash"`
	Duration       float64 `gorose:"duration"`
	Frames         int     `gorose:"frames"`
	Width          int     `gorose:"width"`
	Height         int     `gorose:"height"`
	Status         int     `gorose:"status"`
	UploadTime     string  `gorose:"uploadtime"`
	UploadUid      int     `gorose:"uploaduid"`
	PatientID      string  `gorose:"patientid"`
	MachineID      string  `gorose:"machineid"`
	FolderName     string  `gorose:"foldername"`
	Fcode          string  `gorose:"fcode"`
	IncludeViews   string  `gorose:"includeviews"`
	Keywords       string  `gorose:"keywords"`
	Memo           string  `gorose:"memo"`
	MediaType      string  `gorose:"mediatype"`
	MediaData      string  `gorose:"mediadata"`
	LabelAuthorUid int     `gorose:"labelauthoruid"`
	LabelReviewUid int     `gorose:"labelreviewuid"`
	LabelProgress  int     `gorose:"labelprogress"`
}
type MediaInfoUltrasonicImage struct {
	Width  int
	Height int
}
type MediaInfoUltrasonicVideo struct {
	PathRaw  string
	HashRaw  string
	PathJpgs string
	Encoder  string
}

func initMediaDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"mid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"displayname" TEXT NOT NULL default "",
	"path" TEXT NOT NULL default "",
	"hash" TEXT NOT NULL default "",
	"duration" NUMERIC NOT NULL default 0,
	"frames" INTERGER NOT NULL default 0,
	"width" INTERGER NOT NULL default 0,
	"height" INTERGER NOT NULL default 0,
	"status" INTERGET NOT NULL default 0,
	"uploadtime" TEXT NOT NULL default "2000-01-01 00:00:00",
	"uploaduid" INTERGER NOT NULL default 0,
	"patientid" TEXT NOT NULL default "",
	"machineid" TEXT NOT NULL default "",
	"foldername" TEXT NOT NULL default "",
	"fcode" TEXT NOT NULL default "",
	"includeviews" TEXT NOT NULL default "",
	"keywords" TEXT NOT NULL default "[]",
	"memo" TEXT NOT NULL default "",
	"mediatype" TEXT NOT NULL default "",
	"mediadata" TEXT NOT NULL default "{}",
	"labelprogress" INTERGER NOT NULL default 0,
	"labelauthoruid" INTERGER NOT NULL default 0,
	"labelauthorslid" INTERGER NOT NULL default 0,  
	"labelreviewuid" INTERGER NOT NULL default 0,
	"labelreviewslid" INTERGER NOT NULL default 0)`, global.DefaultDatabaseMediaTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log("e", err.Error())
	}
}

func MediaCreate(mi MediaInfo) (mid int, err error) {
	if mi.Hash == "" {
		return 0, errors.New(global.EMediaRawHashNull)
	}
	mt, err := MediaGet(mi.Hash)
	if err == nil {
		return mt.Mid, errors.New(global.EMediaAlreadyExisted)
	}

	mi.Mid = 0
	// Upgrade override
	// mi.UploadTime = time.Now().Format(global.TimeFormat)
	_, err = DB().Table(global.DefaultDatabaseMediaTable).Data(mi).Insert()

	mi, err = MediaGet(mi.Hash)
	return mi.Mid, err
}

func MediaGet(i interface{}) (mi MediaInfo, err error) {
	switch i.(type) {
	case int:
		err = DB().Table(&mi).Where("mid", "=", i).Select()
		if mi.Mid == 0 {
			if err != nil {
				log("e", err.Error())
			}
			err = errors.New(global.EMediaNotExist)
		}
		return

	case string:
		key := i.(string)
		err = DB().Table(&mi).Where("hash", "=", key).Select()
		if mi.Mid == 0 {
			if err != nil {
				log("E", err.Error())
			}
			err = errors.New(global.EMediaNotExist)
		}
	}
	return
}
func MediaGetByDisplayName(displayname string) (mi MediaInfo, err error) {
	var mis []MediaInfo
	err = DB().Table(&mis).Where("displayname", "=", displayname).Select()

	if err != nil {
		log("E", "MediaGetByDisplayName", err.Error())
	} else if len(mis) == 0 {
		err = errors.New(global.EMediaNotExist)
	} else if len(mis) > 1 {
		err = errors.New("发现媒体使用相同的文件名")
	} else {
		mi = mis[0]
	}
	return
}

func MediaGetAll() (ml []MediaInfo, err error) {
	err = DB().Table(&ml).OrderBy("mid").Select()
	return
}

func MediaUpdate(mid int, data interface{}) (err error) {
	_, err = DB().Table(global.DefaultDatabaseMediaTable).Data(data).Where("mid", "=", mid).Update()
	return
}
func MediaUpdateDetail(mid int, mediaData interface{}) error {
	_, err := MediaGet(mid)
	if err != nil {
		return err
	}

	mediaT := ""
	jb, err := json.Marshal(mediaData)
	if err != nil {
		return err
	}
	mediaD := string(jb)

	switch mediaData.(type) {
	case MediaInfoUltrasonicImage:
		mediaT = global.MediaTypeUltrasonicImage

	case MediaInfoUltrasonicVideo:
		mediaT = global.MediaTypeUltrasonicVideo

	default:
		return errors.New(global.EMediaUnknownType)
	}

	data := map[string]interface{}{"mediatype": mediaT, "mediadata": mediaD}
	return MediaUpdate(mid, data)
}
func MediaUpdatePath(mid int, path string) error {
	data := map[string]interface{}{"path": path}
	return MediaUpdate(mid, data)
}
func MediaUpdateFolderName(mid int, foldername string) error {
	data := map[string]interface{}{"foldername": foldername}
	return MediaUpdate(mid, data)
}
func MediaUpdateDisplayName(mid int, name string) error {
	data := map[string]interface{}{"displayname": name}
	return MediaUpdate(mid, data)
}
func MediaUpdateAuthorReview(mid, authorUid, reviewUid int) error {
	data := map[string]interface{}{"labelauthoruid": authorUid, "labelreviewuid": reviewUid}
	return MediaUpdate(mid, data)
}
func MediaGetDetail(mid int) (mediaData interface{}, err error) {
	mi, err := MediaGet(mid)
	if err != nil {
		return
	}

	switch mi.MediaType {
	case global.MediaTypeUltrasonicVideo:
		var data MediaInfoUltrasonicVideo
		err = json.Unmarshal([]byte(mi.MediaData), &data)
		return data, err

	case global.MediaTypeUltrasonicImage:
		var data MediaInfoUltrasonicImage
		err = json.Unmarshal([]byte(mi.MediaData), &data)
		return data, err

	default:
		return mediaData, errors.New(global.EMediaUnknownType)
	}
}

// Frames and duration
func MediaUpdateFramesAndDuration(mid int, frames int, duration float64) error {
	data := map[string]interface{}{"frames": frames, "duration": duration}
	return MediaUpdate(mid, data)
}
func MediaUpdateWidthAndHeight(mid, width, height int) error {
	data := map[string]interface{}{"width": width, "height": height}
	return MediaUpdate(mid, data)
}

func MediaUpdateHash(mid int, hash string) error {
	data := map[string]interface{}{"hash": hash}
	return MediaUpdate(mid, data)
}

// Views
func MediaGetViews(mid int) (views []string, err error) {
	mi, err := MediaGet(mid)
	if err != nil {
		return
	}
	switch mi.IncludeViews {
	case "", "[]", "null":
		return make([]string, 0), nil
	default:
		views = strings.Split(mi.IncludeViews, ",")
		return
	}
}
func MediaAddView(mid int, view string) error {
	views, err := MediaGetViews(mid)
	if err != nil {
		return err
	}
	for _, v := range views {
		if v == view {
			return nil
		}
	}
	views = append(views, view)
	return MediaSetViews(mid, views)
}
func MediaSetViews(mid int, views []string) error {
	viewStr := ""
	for _, view := range views {
		if viewStr == "" {
			viewStr = view
		} else {
			viewStr = fmt.Sprintf("%s,%s", viewStr, view)
		}
	}
	return MediaUpdateViews(mid, viewStr)
}
func MediaUpdateViews(mid int, viewsStr string) error {
	data := map[string]interface{}{"includeviews": viewsStr}
	return MediaUpdate(mid, data)
}
func MediaRemoveView(mid int, view string) error {
	views, err := MediaGetViews(mid)
	if err != nil {
		return err
	}

	for i, v := range views {
		if v == view {
			views = append(views[:i], views[i+1:]...)
			return MediaSetViews(mid, views)
		}
	}
	return nil
}

// Keywords
func MediaGetKeywords(mid int) (keywords []string, err error) {
	mi, err := MediaGet(mid)
	if err != nil {
		return
	}

	if mi.Keywords == "" || mi.Keywords == "[]" {
		return make([]string, 0), nil
	}

	err = json.Unmarshal([]byte(mi.Keywords), &keywords)
	return
}
func MediaAddKeyword(mid int, keyword string) error {
	keywords, err := MediaGetKeywords(mid)
	if err != nil {
		return err
	}

	for _, v := range keywords {
		if v == keyword {
			return nil
		}
	}

	keywords = append(keywords, keyword)
	return MediaSetKeywords(mid, keywords)
}
func MediaSetKeywords(mid int, keywords []string) error {
	if keywords == nil {
		keywords = make([]string, 0)
	}

	jb, err := json.Marshal(keywords)
	if err != nil {
		return err
	}
	return MediaUpdateKeywords(mid, string(jb))
}
func MediaUpdateKeywords(mid int, keywordStr string) error {
	data := map[string]interface{}{"keywords": keywordStr}
	return MediaUpdate(mid, data)
}
func MediaRemoveKeyword(mid int, keyword string) error {
	keywords, err := MediaGetKeywords(mid)
	if err != nil {
		return err
	}

	for i, v := range keywords {
		if v == keyword {
			keywords = append(keywords[:i], keywords[i+1:]...)
			return MediaSetKeywords(mid, keywords)
		}
	}
	return nil
}

// Memo
func MediaUpdateMemo(mid int, memo string) error {
	data := map[string]interface{}{"memo": memo}
	return MediaUpdate(mid, data)
}

// Label
func MediaUpdateLabelAuthorUidLid(mid int, uid, lid int) error {
	data := map[string]interface{}{"labelauthoruid": uid, "labelauthorslid": lid}
	return MediaUpdate(mid, data)
}
func MediaUpdateLabelReviewUidLid(mid int, uid, lid int) error {
	data := map[string]interface{}{"labelreviewuid": uid, "labelreviewslid": lid}
	return MediaUpdate(mid, data)
}
func MediaUpdateLabelProgress(mid, authoruid, revieweruid, progress int) error {
	data := map[string]interface{}{"labelauthoruid": authoruid, "labelreviewuid": revieweruid, "labelprogress": progress}
	return MediaUpdate(mid, data)
}
func MediaRemoveLabel(mid int) error {
	data := map[string]interface{}{"labelauthoruid": 0, "labelreviewuid": 0, "labelprogress": 0}
	return MediaUpdate(mid, data)
}

func MediaDelete(mid int) (err error) {
	_, err = DB().Table(global.DefaultDatabaseMediaTable).Where("mid", "=", mid).Delete()
	return
}
