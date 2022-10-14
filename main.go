package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"time"
	"uni-minds.com/liuxy/medical-sys/controller"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/module"
)

var _BUILD_TIME_ = "20201022"
var _BUILD_REV_ = "DEBUG"
var _BUILD_VER_ = "2.1"

func main() {
	var argHttps bool
	var argPort int
	var argRegCode string

	config := global.GetAppSettings()

	flag.BoolVar(&argHttps, "s", config.SystemUseHttps, "use https (need certification file)")
	flag.IntVar(&argPort, "p", config.SystemListenPort, "use port")
	flag.StringVar(&argRegCode, "r", config.UserRegisterCode, "register code")
	flag.Parse()

	config.SystemUseHttps = argHttps
	config.SystemListenPort = argPort
	config.UserRegisterCode = argRegCode
	global.SetAppSettings(config)

	t, _ := time.Parse("2006-01-02 15:04:05", _BUILD_TIME_)
	verStr := fmt.Sprintf("%s(%s) %s", _BUILD_VER_, _BUILD_REV_, t.Format("20060102-150405"))
	fmt.Println(color.HiRedString("Version: %s", verStr))
	global.SetVersionString(verStr)

	if err := os.MkdirAll(path.Join(config.SystemAppPath, "log"), os.ModePerm); err != nil {
		panic(err.Error())
	}
	fp1, err := os.OpenFile(path.Join(config.SystemAppPath, "log/access.log"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer fp1.Close()

	fp2, err := os.OpenFile(path.Join(config.SystemAppPath, "log/error.log"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer fp2.Close()

	gin.DefaultWriter = fp1
	gin.DefaultErrorWriter = fp2

	router := gin.Default()
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
		rUi.GET("/home", controller.UIHomeGet)
		rUi.GET("/manage/:class", controller.UIManageGetHandler)
		rUi.GET("/medialist", controller.UIMedialistGet)
		rUi.GET("/labeltool", controller.UILabeltoolGet)
		rUi.GET("/import", controller.UIImportMedia)
	}

	rMobi := router.Group("/mobi", checkUserAuthorized)
	{
		rMobi.GET("/", controller.MobiRoot)
		rMobi.GET("/device", controller.MobiGetDevice)
		rMobi.GET("/result/:pipeline", controller.MobiGetResult)
		rMobi.GET("/exec", controller.MobiMyExec)
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

		apiV1.GET("group", controller.GroupGet)

		apiV1.GET("raw", controller.GetRawData)

		apiV1.GET("blockchain/nodelist", controller.GetBlockchainNodelist)
		apiV1.GET("blockchain/tps", controller.GetBlockchainTPS)
		apiV1.GET("blockchain/height", controller.GetBlockHeight)

		apiV1.Group("database/ct/:class/rs").Any("/*a", controller.GetDatabaseDicomCtRsGroup)
		apiV1.Group("database/ct/:class/wado").Any("/*a", controller.GetDatabaseDicomCtWadoGroup)
		apiV1.POST("analysis/ct/:class/:mode", controller.AnalysisCtPost)

	}

	certPem := path.Join(config.SystemAppPath, "server.crt")
	certKey := path.Join(config.SystemAppPath, "server.key")
	strPort := fmt.Sprintf(":%d", config.SystemListenPort)

	if argHttps {
		log.Println(color.RedString("Use HTTPS @ Port %s", strPort))
		router.RunTLS(strPort, certPem, certKey)
	} else {
		log.Println(color.RedString("Use HTTP @ Port %s", strPort))
		router.Run(strPort)
	}
}

func checkUserAuthorized(ctx *gin.Context) {
	valid, _ := controller.CookieValidUid(ctx)
	if !valid {
		log.Println(color.HiRedString("Unauthorized"))
		ctx.Redirect(http.StatusFound, "/")
	}
}
