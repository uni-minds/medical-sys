/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: init.go
 */

package database

import (
	"fmt"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
	"path/filepath"
	"uni-minds.com/liuxy/medical-sys/global"
)

var engin *gorose.Engin

func init() {
	var err error

	c := global.GetAppSettings()
	dbfile, _ := filepath.Abs(c.SystemDBFile)
	fmt.Println("DB:", dbfile)
	if _, err := os.Stat(path.Dir(dbfile)); err != nil {
		err = os.MkdirAll(path.Dir(dbfile), os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}

	engin, err = gorose.Open(&gorose.Config{
		Driver: "sqlite3",
		Dsn:    dbfile,
	})

	if err != nil {
		log.Panic(err.Error())
	}

	initUserDB()
	initMediaDB()
	initGroupDB()
	initLabelDB()
}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
