/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: class_label.go
 * Date: 2022/09/23 21:32:23
 */

package database

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/utils/tools"
	"os"
	"runtime"
)

var dbLabel *DbLabel

type DbLabel struct {
	TableName  string
	ModuleName string
}

func GetLabel() *DbLabel {
	if dbLabel == nil {
		dbLabel = CreateLabelDB(os.Stdout)
	}
	return dbLabel
}

func CreateLabelDB(fp *os.File) *DbLabel {
	db := new(DbLabel)
	db.TableName = global.DefaultDatabaseLabelTable
	db.ModuleName = "DB_L"
	return db
}

func InitLabelDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
    id                 INTEGER            not null
        constraint labels_pk
            primary key autoincrement,
    progress           INT     default 0  not null,
    author             INT     default 0  not null,
    reviewer           INT     default 0  not null,
    media_uuid         TEXT               not null,
    data               TEXT               not null,
    version            INT     default 1  not null,
    frames             INTEGER default 0  not null,
    counts             INTEGER default 0  not null,
    t_author_start     TEXT    default "" not null,
    t_author_submit    TEXT    default "" not null,
    t_reviewer_start   TEXT    default "" not null,
    t_reviewer_confirm TEXT    default "" not null,
    memo               TEXT    default "" not null,
    t_reviewer_submit  TEXT    default "" not null,
    label_uuid         TEXT    default '' not null)`, global.DefaultDatabaseLabelTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Error(err.Error())
	}
}

func (p *DbLabel) Trace(msg ...any) {
	p.Write("t", tools.InterfaceExpand(msg))
}
func (p *DbLabel) Warn(msg ...any) {
	p.Write("w", tools.InterfaceExpand(msg))
}
func (p *DbLabel) Error(msg ...any) {
	p.Write("e", tools.InterfaceExpand(msg))
}
func (p *DbLabel) Write(level string, msg string) {
	_, file, line, ok := runtime.Caller(2)
	if ok && global.FlagGetDebug() {
		logger.Write(p.ModuleName, level, fmt.Sprintf("%s:%d %s", file, line, msg))
	} else {
		logger.Write(p.ModuleName, level, msg)
	}
}

func (p *DbLabel) Get(i interface{}) (info DbStructLabel, err error) {
	switch i.(type) {
	case int:
		if err = DB().Table(&info).Where("id", i).Select(); err != nil {
			return info, err
		}

	case string:
		if err = DB().Table(&info).Where("label_uuid", i).Select(); err != nil {
			return info, err
		}
	}

	if info.Id < 1 {
		return info, fmt.Errorf(global.ELabelDBLabedNotExist)
	} else {
		return info, nil
	}
}
func (p *DbLabel) GetByMediaUUID(mediaUUID string) (info []DbStructLabel, err error) {
	err = DB().Table(&info).Where("media_uuid", mediaUUID).Select()
	return info, err
}
func (p *DbLabel) GetAll() (infos []DbStructLabel, err error) {
	err = DB().Table(&infos).OrderBy("id").Select()
	return infos, err
}

func (p *DbLabel) Create(li DbStructLabel) error {
	p.Trace("create:", li)
	li.Id = 0
	_, err := DB().Table(p.TableName).Data(li).Insert()
	if err != nil {
		p.Error("Label create:", err.Error())
	}
	return err
}

func (p *DbLabel) UpdateAll(li DbStructLabel) error {
	_, err := DB().Table(p.TableName).Data(li).Where("id", li.Id).Update()
	return err
}
func (p *DbLabel) UpdateManual(i interface{}, data map[string]interface{}) error {
	//data["timeAuthorSubmit"] = time.Now().Format(global.TimeFormat)
	info, err := p.Get(i)
	if err != nil {
		return err
	}

	_, err = DB().Table(p.TableName).Data(data).Where("id", info.Id).Update()
	return err
}

func (p *DbLabel) UpdateMemo(i interface{}, memo string) error {
	return p.UpdateManual(i, map[string]interface{}{"memo": memo})
}
func (p *DbLabel) UpdateJsonData(i interface{}, jsonstr string) error {
	data := map[string]interface{}{"data": jsonstr}
	return p.UpdateManual(i, data)
}
func (p *DbLabel) UpdateProgress(i interface{}, progress int) error {
	data := map[string]interface{}{"progress": progress}
	return p.UpdateManual(i, data)
}

func (p *DbLabel) Delete(i interface{}) (err error) {
	info, err := p.Get(i)
	if err != nil {
		return err
	}
	_, err = DB().Table(p.TableName).Where("id", info.Id).Delete()
	return err
}
