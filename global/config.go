package global

import (
	"fmt"
	"time"
	"uni-minds.com/liuxy/medical-sys/tools"
)

var config AppSettings
var configFile = "config.yaml"

type AppSettings struct {
	UserRegisterEnable bool
	UserRegisterCode   string
	CookieMaxAge       int
	SystemListenPort   int
	SystemUseHttps     bool
	SystemAppPath      string
	SystemMediaPath    string
	SystemDBFile       string
	SystemLogFolder    string
}

type Version struct {
	GitCommit string
	Version   string
	BuildTime string
}

func init() {
	config = AppSettings{
		UserRegisterEnable: true,
		UserRegisterCode:   "beihang",
		CookieMaxAge:       24 * int(time.Hour.Seconds()),
		SystemListenPort:   80,
		SystemUseHttps:     false,
		SystemAppPath:      "/usr/local/uni-ledger/medical-sys/application",
		SystemMediaPath:    "/usr/local/uni-ledger/medical-sys/application/media",
		SystemDBFile:       "/usr/local/uni-ledger/medical-sys/application/database/db.sqlite",
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
