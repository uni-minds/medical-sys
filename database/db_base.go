/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: init.go
 */

package database

import (
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

var engin *gorose.Engin
var log *logger.Logger

func Init() (err error) {
	log = logger.NewLogger("DB")

	dbfile, _ := global.GetDbFile("main")
	log.Println("main db ->", dbfile)
	if _, err = os.Stat(path.Dir(dbfile)); err != nil {
		err = os.MkdirAll(path.Dir(dbfile), os.ModePerm)
		if err != nil {
			return err
		}
	}

	engin, err = gorose.Open(&gorose.Config{
		Driver: "sqlite3",
		Dsn:    dbfile,
	})

	if err != nil {
		return err
	}

	initUserDB()
	InitMediaDB()
	initGroupDB()
	InitLabelDB()

	BridgeHisInit()
	BridgePacsInit()

	return nil
}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
