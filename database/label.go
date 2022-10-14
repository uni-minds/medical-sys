package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"uni-minds.com/medical-sys/global"
	"uni-minds.com/medical-sys/tools"
)

func (*LabelInfo) TableName() string {
	return global.DefaultDatabaseLabelTable
}

type LabelInfo struct {
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

func initLabelDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"lid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"uid" INTERGER NOT NULL default 0,
	"mid" INTERGER NOT NULL default 0,
	"type" TEXT NOT NULL default "",
	"data" TEXT NOT NULL default "",
	"databackup" TEXT NOT NULL default "",
	"version" INTERGER NOT NULL default 0,
	"progress" TEXT NOT NULL default "",
	"frames" INTERGER NOT NULL default 0,
	"counts" INTERGER NOT NULL default 0,
	"createtime" TEXT NOT NULL default "",
	"modifytime" TEXT NOT NULL default "",
	"memo" TEXT NOT NULL default "",
	"hash" TEXT NOT NULL UNIQUE)`, global.DefaultDatabaseLabelTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Panic(err.Error())
	}
}
func LabelCreate(li LabelInfo) (lid int, labelUUID string, err error) {
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

func LabelQuery(mid, uid int, ltype string) (li LabelInfo, err error) {
	err = DB().Table(&li).Where("uid", "=", uid).Where("mid", "=", mid).Where("type", "=", ltype).Select()
	if err != nil || li.Lid == 0 {
		err = errors.New(global.ELabelDBLabedNotExist)
	}
	return
}
func LabelGet(i interface{}) (li LabelInfo, err error) {
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
func LabelGetAll(mid, uid int, labeltype string) (lis []LabelInfo, err error) {
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

func labelUpdate(lid int, data interface{}) (err error) {
	d := data.(map[string]interface{})
	d["modifytime"] = time.Now().Format(global.TimeFormat)
	_, err = DB().Table(global.DefaultDatabaseLabelTable).Data(d).Where("lid", "=", lid).Update()
	if err != nil {
		fmt.Println("DB E", err.Error())
	}
	return
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
func LabelUpdateMemo(lid int, memo string) error {
	data := map[string]interface{}{"memo": memo}
	return labelUpdate(lid, data)
}

func LabelDelete(lid int) (err error) {
	_, err = DB().Table(global.DefaultDatabaseLabelTable).Where("lid", "=", lid).Delete()
	return
}
