/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: config.go
 */

package global

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/tools"
	"time"
)

var config AppSettings
var configFile = "/usr/local/uni-ledger/medical-sys/config.yaml"

type AppSettings struct {
	UserRegisterEnable bool
	UserRegisterCode   string
	CookieMaxAge       int
	SystemListenPort   int
	SystemUseHttps     bool
	PathApp            string
	PathMedia          string
	PathPacsMediaCache string
	DbFileMain         string
	DbFilePacs         string
	DbFileHis          string
	SystemLogFolder    string
}

type Version struct {
	GitCommit string
	Version   string
	BuildTime string
}

func init() {
	fmt.Println("module init: global")
	config = AppSettings{
		UserRegisterEnable: true,
		UserRegisterCode:   "BUAA",
		CookieMaxAge:       24 * int(time.Hour.Seconds()),
		SystemListenPort:   80,
		SystemUseHttps:     false,
		PathApp:            "/usr/local/uni-ledger/medical-sys/application",
		PathMedia:          "/usr/local/uni-ledger/medical-sys/application/media",
		PathPacsMediaCache: "/data/cache",
		DbFileMain:         "/usr/local/uni-ledger/medical-sys/application/database/db.sqlite",
		DbFilePacs:         "/usr/local/uni-ledger/medical-sys/application/database/db_pacs.sqlite",
		DbFileHis:          "/usr/local/uni-ledger/medical-sys/application/database/db_his.sqlite",
		SystemLogFolder:    "/usr/local/uni-ledger/medical-sys/log",
	}
	loadConfig(configFile)
}

func loadConfig(file string) {
	var c AppSettings
	if err := tools.LoadYaml(file, &c); err == nil {
		config = c
	} else {
		fmt.Println("E;Load CFG:", err.Error())
		saveConfig(file)
	}
}

func saveConfig(file string) {
	tools.SaveYaml(file, config)
	configFile = file
}

func GetAppSettings() AppSettings {
	return config
}

func SetAppSettings(data AppSettings) {
	config = data
	saveConfig(configFile)
}

func GetCookieMaxAge() int {
	return config.CookieMaxAge
}

func GetUserRegCode() string {
	return config.UserRegisterCode
}
