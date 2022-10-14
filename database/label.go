/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: label.go
 */

/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  labels
 * @Version: 1.0.0
 * @Date: 2020/4/14 23:50
 */

package database

import (
	"errors"
	"fmt"
	"log"
	"time"
	"uni-minds.com/liuxy/medical-sys/global"
)

func (*LabelInfo) TableName() string {
	return global.DefaultDatabaseLabelTable
}

type LabelInfo struct {
	Lid               int    `gorose:"lid"`
	Progress          int    `gorose:"progress"`
	AuthorUid         int    `gorose:"authorUid"`
	ReviewUid         int    `gorose:"reviewUid"`
	MediaHash         string `gorose:"mediaHash"`
	Data              string `gorose:"data"`
	Version           int    `gorose:"version"`
	Frames            int    `gorose:"frames"`
	Counts            int    `gorose:"counts"`
	TimeAuthorStart   string `gorose:"timeAuthorStart"`
	TimeAuthorSubmit  string `gorose:"timeAuthorSubmit"`
	TimeReviewStart   string `gorose:"timeReviewStart"`
	TimeReviewSubmit  string `gorose:"timeReviewSubmit"`
	TimeReviewConfirm string `gorose:"timeReviewConfirm"`
	Memo              string `gorose:"memo"`
}

func initLabelDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"lid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"progress" TEXT NOT NULL default "",
	"authorUid" INTERGER NOT NULL default 0,
	"reviewUid" INTERGER NOT NULL default 0,
	"mediaHash" TEXT NOT NULL UNIQUE,
	"data" TEXT NOT NULL default "",
	"version" INTERGER NOT NULL default 0,
	"frames" INTERGER NOT NULL default 0,
	"counts" INTERGER NOT NULL default 0,
	"timeAuthorStart" TEXT NOT NULL default "",
	"timeAuthorSubmit" TEXT NOT NULL default "",
	"timeReviewStart" TEXT NOT NULL default "",
	"timeReviewSubmit" TEXT NOT NULL default "",
	"timeReviewConfirm" TEXT NOT NULL default "",
	"memo" TEXT NOT NULL default "")`, global.DefaultDatabaseLabelTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Panic(err.Error())
	}
}
func LabelGet(i interface{}) (li LabelInfo, err error) {
	switch i.(type) {
	case int:
		err = DB().Table(&li).Where("lid", "=", i).Select()
		if err != nil || li.Lid == 0 {
			err = errors.New(global.ELabelDBLabedNotExist)
		}
	case string:
		err = DB().Table(&li).Where("mediahash", "=", i).Select()
		if err != nil || li.Lid == 0 {
			err = errors.New(global.ELabelDBLabedNotExist)
		}
	}
	return
}
func LabelUpdateProgress(lid, progress int) error {
	data := map[string]interface{}{"progress": progress}
	_, err := DB().Table(global.DefaultDatabaseLabelTable).Data(data).Where("lid", "=", lid).Update()
	return err
}
func LabelCreate(li LabelInfo) error {
	li.Lid = 0
	_, err := DB().Table(global.DefaultDatabaseLabelTable).Data(li).Insert()
	if err != nil {
		log.Println("E Label create:", err.Error())
	}
	return err
}
func LabelUpdate(li LabelInfo) (err error) {
	_, err = DB().Table(global.DefaultDatabaseLabelTable).Data(li).Where("lid", "=", li.Lid).Update()
	if err != nil {
		fmt.Println("DB E", err.Error())
	}
	return
}
func labelUpdate(lid int, data interface{}) (err error) {
	d := data.(map[string]interface{})
	d["timeAuthorSubmit"] = time.Now().Format(global.TimeFormat)
	_, err = DB().Table(global.DefaultDatabaseLabelTable).Data(d).Where("lid", "=", lid).Update()
	if err != nil {
		fmt.Println("DB E", err.Error())
	}
	return
}
func LabelUpdateMemo(lid int, memo string) error {
	data := map[string]interface{}{"memo": memo}
	return labelUpdate(lid, data)
}
func LabelUpdateMediaHash(lid int, mediaHash string) error {
	data := map[string]interface{}{"mediaHash": mediaHash}
	return labelUpdate(lid, data)
}
func LabelUpdateJsonDataOnly(lid int, jsonstr string) error {
	data := map[string]interface{}{"data": jsonstr}
	_, err := DB().Table(global.DefaultDatabaseLabelTable).Data(data).Where("lid", "=", lid).Update()
	return err
}

func LabelQuery(mid, uid int, ltype string) (li LabelInfo, err error) {
	err = DB().Table(&li).Where("uid", "=", uid).Where("mid", "=", mid).Where("type", "=", ltype).Select()
	if err != nil || li.Lid == 0 {
		err = errors.New(global.ELabelDBLabedNotExist)
	}
	return
}
func LabelGetAll() (lis []LabelInfo, err error) {
	err = DB().Table(&lis).OrderBy("lid").Select()
	return
}
func LabelDelete(i interface{}) (err error) {
	switch i.(type) {
	case int:
		_, err = DB().Table(global.DefaultDatabaseLabelTable).Where("lid", "=", i.(int)).Delete()
	case string:
		_, err = DB().Table(global.DefaultDatabaseLabelTable).Where("mediaHash", "=", i.(string)).Delete()
	}
	return err
}

/*
type _LabelInfo struct {
	Lid        int    `gorose:"lid"`
	Uid        int    `gorose:"uid"`
	Mid        int    `gorose:"mid"`
	Type       string `gorose:"type"`
	Data       string `gorose:"data"`
	DataBackup string `gorose:"databackup"`
	Version    int    `gorose:"version"`
	Progress   string `gorose:"progress"`
	Frames     int    `gorose:"frames"`
	Counts     int    `gorose:"counts"`
	CreateTime string `gorose:"createtime"`
	ModifyTime string `gorose:"modifytime"`
	Memo       string `gorose:"memo"`
	Hash       string `gorose:"hash"`
}
func LabelUpdateLabelData(lid, uid, frames, counts int, labeldata interface{}, progress string) error {
	li, err := LabelGet(lid)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	switch labeldata.(type) {
	case LabelInfoAuthorData:
		ldata, _ := json.Marshal(labeldata)
		dataStr := string(ldata)
		//hash := tools.GetStringMD5(dataStr + strconv.Itoa(lid))
		data = map[string]interface{}{"data": dataStr, "databackup": li.Data, "type": global.LabelTypeAuthor,
			"uid": uid, "progress": progress, "frames": frames, "counts": counts}
		//, "hash": hash}
	case LabelInfoReviewerData:
		ldata, _ := json.Marshal(labeldata)
		dataStr := string(ldata)
		//hash := tools.GetStringMD5(dataStr + strconv.Itoa(lid))
		data = map[string]interface{}{"data": dataStr, "databackup": li.Data, "type": global.LabelTypeReview,
			"uid": uid, "progress": progress, "frames": frames, "counts": counts}
		//, "hash": hash}
	default:
		return errors.New("not a regular label type.")
	}

	return labelUpdate(lid, data)
}

type LabelInfoAuthorData struct {
	Json              string
	Reviewed          bool
	ModifyAfterReview bool
}
type LabelInfoReviewerData struct {
	BasedJson   string
	BasedAuthor int
	BasedTime   string
	Json        string
}
func _LabelCreate(li LabelInfo) (lid int, labelUUID string, err error) {
	switch li.Type {
	case global.LabelTypeAuthor, global.LabelTypeReview, global.LabelTypeFinal:
		_, err = LabelQuery(li.Mid, li.Uid, li.Type)
		if err == nil {
			return 0, "", errors.New(global.ELabelDBLabelExisted)
		}
		li.Lid = 0
		li.CreateTime = time.Now().Format(global.TimeFormat)
		log.Println("随机生成HASH作为Label UUID")
		li.Hash = tools.GetStringMD5(tools.GenSaltString(20))
		_, err = DB().Table(global.DefaultDatabaseLabelTable).Data(li).Insert()
		if err != nil {
			return 0, "", err
		}

		li, err = LabelQuery(li.Mid, li.Uid, li.Type)
		return li.Lid, li.Hash, err
	default:
		return 0, "", errors.New(global.ELabelInvalidType)
	}
}

func _LabelQuery(mid, uid int, ltype string) (li LabelInfo, err error) {
	err = DB().Table(&li).Where("uid", "=", uid).Where("mid", "=", mid).Where("type", "=", ltype).Select()
	if err != nil || li.Lid == 0 {
		err = errors.New(global.ELabelDBLabedNotExist)
	}
	return
}
func __LabelGet(i interface{}) (li LabelInfo, err error) {
	switch i.(type) {
	case int:
		err = DB().Table(&li).Where("lid", "=", i).Select()
		if err != nil || li.Lid == 0 {
			err = errors.New(global.ELabelDBLabedNotExist)
		}
	case string:
		err = DB().Table(&li).Where("hash", "=", i).Select()
		if err != nil || li.Lid == 0 {
			err = errors.New(global.ELabelDBLabedNotExist)
		}
	}
	return
}
func _LabelGetAll(mid, uid int, labeltype string) (lis []LabelInfo, err error) {
	iOrm := DB().Table(&lis)
	if mid > 0 {
		iOrm.Where("mid", "=", mid)
	}

	if uid > 0 {
		iOrm.Where("uid", "=", uid)
	}

	if labeltype != "" {
		iOrm.Where("type", "=", labeltype)
	}

	err = iOrm.OrderBy("lid").Select()
	return
}

func _labelUpdate(lid int, data interface{}) (err error) {
	d := data.(map[string]interface{})
	d["modifytime"] = time.Now().Format(global.TimeFormat)
	_, err = DB().Table(global.DefaultDatabaseLabelTable).Data(d).Where("lid", "=", lid).Update()
	if err != nil {
		fmt.Println("DB E", err.Error())
	}
	return
}
func _LabelUpdateLabelData(lid, uid, frames, counts int, labeldata interface{}, progress string) error {
	li, err := LabelGet(lid)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	switch labeldata.(type) {
	case LabelInfoAuthorData:
		ldata, _ := json.Marshal(labeldata)
		dataStr := string(ldata)
		//hash := tools.GetStringMD5(dataStr + strconv.Itoa(lid))
		data = map[string]interface{}{"data": dataStr, "databackup": li.Data, "type": global.LabelTypeAuthor,
			"uid": uid, "progress": progress, "frames": frames, "counts": counts}
		//, "hash": hash}
	case LabelInfoReviewerData:
		ldata, _ := json.Marshal(labeldata)
		dataStr := string(ldata)
		//hash := tools.GetStringMD5(dataStr + strconv.Itoa(lid))
		data = map[string]interface{}{"data": dataStr, "databackup": li.Data, "type": global.LabelTypeReview,
			"uid": uid, "progress": progress, "frames": frames, "counts": counts}
		//, "hash": hash}
	default:
		return errors.New("not a regular label type.")
	}

	return labelUpdate(lid, data)
}
func _LabelUpdateMemo(lid int, memo string) error {
	data := map[string]interface{}{"memo": memo}
	return labelUpdate(lid, data)
}

func LabelDelete(lid int) (err error) {
	_, err = DB().Table(global.DefaultDatabaseLabelTable).Where("lid", "=", lid).Delete()
	return
}


*/
