/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: main.go
 */

package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"path"
	"time"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/logger"
	"uni-minds.com/liuxy/medical-sys/main_core"
)

var (
	_BUILD_VER_  = "0.0.0"
	_BUILD_TIME_ = ""
	_BUILD_REV_  = "DEV"
)

func main() {
	var argHttps, argVerbose, argDebug bool
	var argPort int
	var argRegCode string

	config := global.GetAppSettings()

	flag.BoolVar(&argDebug, "debug", false, "Set debug mode (golden token enable)")
	flag.BoolVar(&argVerbose, "v", false, "Verbose")
	flag.BoolVar(&argHttps, "s", config.SystemUseHttps, "use https (need certs file)")
	flag.IntVar(&argPort, "p", config.SystemListenPort, "use port")
	flag.StringVar(&argRegCode, "r", config.UserRegisterCode, "register code")
	flag.Parse()

	config.SystemUseHttps = argHttps
	config.SystemListenPort = argPort
	config.UserRegisterCode = argRegCode
	global.SetAppSettings(config)

	t, _ := time.Parse("2006-01-02 15:04:05", _BUILD_TIME_)
	version := fmt.Sprintf("%s.%s_arm64_%s", _BUILD_VER_, t.Format("20060102"), _BUILD_REV_)
	global.SetVersionString(version)
	logo(version)

	logger.Init(path.Join(config.SystemLogFolder, "medical-sys.log"), argVerbose)
	global.DebugSetFlag(argDebug)

	main_core.Router()
}

func logo(version string) {
	fmt.Println(color.HiGreenString("    __  ___         ___            __                     \n   /  |/  /__  ____/ (_)________ _/ /     _______  _______\n  / /|_/ / _ \\/ __  / / ___/ __ `/ /_____/ ___/ / / / ___/\n / /  / /  __/ /_/ / / /__/ /_/ / /_____(__  ) /_/ (__  ) \n/_/  /_/\\___/\\__,_/_/\\___/\\__,_/_/     /____/\\__, /____/  \n                                            /____/"))
	fmt.Println(color.HiMagentaString("Beihang university medical-sys %s", version))
}
