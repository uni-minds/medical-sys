/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: main.go
 * Description:
 */

package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"time"
	"uni-minds.com/liuxy/medical-sys/controller"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/logger"
	"uni-minds.com/liuxy/medical-sys/module"
	"uni-minds.com/liuxy/medical-sys/tools"
)

const tag = "MAIN"

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

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	//router.Use(timeoutMiddleware(5 * time.Second))
	router.Use(gin.Recovery(), gin.ErrorLogger())
	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		switch params.StatusCode {
		case 200:
			logger.Write("HTTP", "t", fmt.Sprintf("%-4s 200 %s", params.Method, params.Path))
		default:
			logger.Write("HTTP", "w", fmt.Sprintf("%-4s %d %s", params.Method, params.StatusCode, params.Path))
		}
		return ""
	}))

	module.Init()
	{
		router.Static("/build", path.Join(config.SystemAppPath, "web/build"))
		router.Static("/pages", path.Join(config.SystemAppPath, "web/pages"))
		router.Static("/application", config.SystemAppPath)
		router.Static("/plugins", path.Join(config.SystemAppPath, "web/plugins"))
		router.Static("/dist", path.Join(config.SystemAppPath, "web/dist"))
		router.Static("/webapp", path.Join(config.SystemAppPath, "web/webapp"))
		router.StaticFile("/favicon.ico", path.Join(config.SystemAppPath, "web/favicon.ico"))
		router.LoadHTMLGlob(path.Join(config.SystemAppPath, "web/templates/**/*"))

		router.GET("/", controller.RootGetHandler)
		router.GET("/login", controller.RootUserLoginGet)
		router.GET("/logout", controller.RootUserLogoutGet)
		router.GET("/register", controller.RootUserRegisterGet)

		router.POST("/register", controller.RootUserRegisterPost)
	}

	rUi := router.Group("/ui", checkUserAuthorized)
	{
		rUi.GET("/home", controller.UiHomeGet)
		rUi.GET("/manage/:class", controller.UiManageGetHandler)
		rUi.GET("/medialist", controller.UiMedialistGet)
		rUi.GET("/labeltool", controller.UiLabeltoolGet)
		rUi.GET("/import", controller.UiImportMedia)
		rUi.GET("/analysis", controller.UiAnalysisGet)
	}

	rMobi := router.Group("/mobi", checkUserAuthorized)
	{
		rMobi.GET("/", controller.MobiRoot)
		rMobi.GET("/device", controller.MobiGetDevice)
		rMobi.GET("/result/:pipeline", controller.MobiGetResult)
		rMobi.POST("/exec", controller.MobiMyExec)
	}

	apiV1 := router.Group("/api/v1", checkUserAuthorized)
	{
		apiV1.POST("login", controller.LoginPost) // Login
		apiV1.GET("login", controller.LoginGet)

		apiV1.POST("user", controller.UserPost)     // CREATE USER
		apiV1.GET("user", controller.UserGet)       // READ USER
		apiV1.DELETE("user", controller.UserDelete) // REMOVE USER
		apiV1.PUT("user", controller.UserPut)       // UPDATE USER

		apiV1.GET("media", controller.MediaGet)

		apiV1.GET("label", controller.LabelGet)
		apiV1.POST("label", controller.LabelPost)
		apiV1.DELETE("label", controller.LabelDel)

		apiV1.GET("algo", controller.AlgoGet)
		apiV1.POST("algo", controller.AlgoPost)

		apiV1.GET("group", controller.GroupGet)

		apiV1.GET("raw", controller.RawDataGet)

		apiV1.GET("blockchain/nodelist", controller.GetBlockchainNodelist)
		apiV1.GET("blockchain/tps", controller.GetBlockchainTPS)
		apiV1.GET("blockchain/height", controller.GetBlockHeight)

		pacs := apiV1.Group("pacs")
		{
			pacs.Any("/", controller.PacsGet)
			pacs.Any("/:db", controller.PacsGetNodelist)
			pacs.Any("/:db/:node", controller.PacsSearchProxy)
			pacs.Any("/:db/:node/:studyuid", controller.PacsSearchProxy)
			pacs.Any("/:db/:node/:studyuid/:seriesuid", controller.PacsSearchProxy)
			pacs.Any("/:db/:node/:studyuid/:seriesuid/:objectuid", controller.PacsSearchProxy)
			//pacs.Any("/:db/:node/wado", controller.PacsWado)
		}

		apiV1.POST("analysis/ct/:class/:mode", controller.AnalysisCtPost)

		aiApi := apiV1.Group("ai")
		aiApi.GET(":modal/:class/:algo/:aid", controller.AiAlgoGet)
		aiApi.GET(":modal/:class/:algo/:aid/:part", controller.AiAlgoGet)
		aiApi.POST(":modal/:class/:algo", controller.AiAlgoPost)
	}

	router.GET("/ws", controller.WebSocket)

	strPort := fmt.Sprintf(":%d", config.SystemListenPort)

	if argHttps {
		log("i", "use https @ port", strPort)
		certPem := path.Join(config.SystemAppPath, "certs/server.crt")
		certKey := path.Join(config.SystemAppPath, "certs/server.key")
		if _, err := os.Stat(certKey); err != nil {
			panic("certs not found")
		} else {
			if err = router.RunTLS(strPort, certPem, certKey); err != nil {
				log("e", err.Error())
			}
		}
	} else {
		log("w", "use http @ port", strPort)
		if err := router.Run(strPort); err != nil {
			log("e", err.Error())
		}
	}
}

func checkUserAuthorized(ctx *gin.Context) {
	key := ctx.Query("key")
	if global.DebugGetFlag() && key == "ADMIN_BUAA" {
		log("w", "Debug ignore user check")
		ctx.Next()
	} else if valid, _ := controller.CookieValidUid(ctx); !valid {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func logo(version string) {
	fmt.Println(color.HiGreenString("    __  ___         ___            __                     \n   /  |/  /__  ____/ (_)________ _/ /     _______  _______\n  / /|_/ / _ \\/ __  / / ___/ __ `/ /_____/ ___/ / / / ___/\n / /  / /  __/ /_/ / / /__/ /_/ / /_____(__  ) /_/ (__  ) \n/_/  /_/\\___/\\__,_/_/\\___/\\__,_/_/     /____/\\__, /____/  \n                                            /____/"))
	fmt.Println(color.HiMagentaString("Beihang university medical-sys %s", version))
}

func log(level string, message ...interface{}) {
	msg := tools.ExpandInterface(message)
	logger.Write(tag, level, msg)
}
