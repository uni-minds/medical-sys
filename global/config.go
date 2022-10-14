/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: config.go
 */

package global

import (
	"fmt"
	utils_tools "gitee.com/uni-minds/utils/tools"
	"path"
	"path/filepath"
)

var config AppSettings
var configFile = "/usr/local/uni-ledger/medical-sys/config.yaml"

type Rtsp struct {
	UploadRoot           string
	NetworkBuffer        int
	Timeout              int
	CloseOld             bool
	PlayerQueueLimit     int
	DropPacketWhenPaused bool
	GopCacheEnable       bool
	DebugLogEnable       bool
	SaveStreamToLocal    bool
	FFmpegPath           string
	FFmpegEncoder        string
	TsDurationSecond     int
	AuthorizationEnable  bool
	ClientUser           string
	ClientPassword       string
}

type UserRegister struct {
	Enable  bool
	RegCode string
}

type Ports struct {
	HTTP int
	RTSP int
	RPC  int
}

type Paths struct {
	Application string
	Media       string
	Cache       string
	Log         string
	Database    string
}

type AppSettings struct {
	CookieMaxAge int
	UserRegister
	Ports
	Rtsp
	Paths
}

func Init(CfgFile string) AppSettings {
	cfg, _ := filepath.Abs(CfgFile)
	fmt.Printf("CFG File: %s\n", cfg)

	config = getDefaultConfig()
	return loadConfig(CfgFile)
}

func loadConfig(file string) AppSettings {
	var c AppSettings
	if err := utils_tools.LoadYaml(file, &c); err == nil {
		config = c
	} else {
		fmt.Println("E;Load CFG:", err.Error())
	}
	saveConfig(file)
	return c
}

func saveConfig(file string) {
	utils_tools.SaveYaml(file, config)
	configFile = file
}

func GetAppSettings() AppSettings {
	return config
}

func GetPaths() Paths {
	return config.Paths
}

func GetDbFile(selector string) (string, error) {
	switch selector {
	case "main":
		return filepath.Abs(path.Join(config.Paths.Database, "db.sqlite"))
	case "pacs":
		return filepath.Abs(path.Join(config.Paths.Database, "db_pacs.sqlite"))
	case "his":
		return filepath.Abs(path.Join(config.Paths.Database, "db_his.sqlite"))
	default:
		return "", fmt.Errorf("unknown db selector")
	}
}

func GetRtspSettings() Rtsp {
	return config.Rtsp
}

func SetAppSettings(data AppSettings) {
	config = data
	saveConfig(configFile)
}

func GetCookieMaxAge() int {
	return config.CookieMaxAge
}

func GetUserRegCode() string {
	return config.UserRegister.RegCode
}
