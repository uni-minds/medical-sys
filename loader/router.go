package main

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/controller"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func router() {
	config := global.GetAppSettings()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
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

	router.Static("/build", path.Join(config.PathApp, "web/build"))
	router.Static("/pages", path.Join(config.PathApp, "web/pages"))
	router.Static("/application", config.PathApp)
	router.Static("/plugins", path.Join(config.PathApp, "web/plugins"))
	router.Static("/plugins-custom", path.Join(config.PathApp, "web/plugins-custom"))
	router.Static("/dist", path.Join(config.PathApp, "web/dist"))
	router.Static("/webapp", path.Join(config.PathApp, "web/webapp"))
	router.StaticFile("/favicon.ico", path.Join(config.PathApp, "web/favicon.ico"))
	router.LoadHTMLGlob(path.Join(config.PathApp, "web/templates/**/*"))

	router.GET("/", controller.RootGet)
	// login logout forget register
	router.GET("/:op", controller.RootGet)
	// user api
	router.GET("/api/user/:op", controller.ApiUserGet)
	router.POST("/api/user/:op", controller.ApiUserPost)

	rUi := router.Group("/ui", TokenCheck)
	{
		rUi.GET("/home", controller.UiHomeGet)
		rUi.GET("/manage/:class", controller.UiManageGetHandler)
		rUi.GET("/medialist", controller.UiMediaListGet)
		rUi.GET("/screen", controller.UiMediaScreenGet)
		rUi.GET("/labeltool/media/:mediaIndex/:usertype", controller.UiLabeltoolGet)
		rUi.GET("/import", controller.UiImportMedia)
		rUi.GET("/analysis", controller.UiAnalysisGet)
	}

	rMobi := router.Group("/mobi", TokenCheck)
	{
		rMobi.GET("/", controller.MobiRoot)
		rMobi.GET("/device", controller.MobiGetDevice)
		rMobi.GET("/result/:pipeline", controller.MobiGetResult)
		rMobi.POST("/exec", controller.MobiMyExec)
	}

	apiV1 := router.Group("/api/v1", TokenCheck)
	{
		apiV1.GET("user", controller.UserGet) // READ USER

		apiV1.GET("media", controller.MediaGet)

		// /api/v1/media/1.1.2.3.4/label | lock |/memo

		apiMedia := apiV1.Group("/media/:mediaIndex", controller.MediaPreOperation)

		apiMedia.GET(":op", controller.MediaGetOperation)
		apiMedia.POST(":op", controller.MediaPostOperation)
		apiMedia.DELETE(":op", controller.MediaDeleteOperation)

		apiMedia.GET("label/:op", controller.LabelGet)
		apiMedia.POST("label/:op", controller.LabelPost)
		apiMedia.DELETE("label/:op", controller.LabelDelete)

		apiV1.GET("screen", controller.ScreenGet)
		apiV1.POST("screen", controller.ScreenPost)
		apiV1.DELETE("screen", controller.ScreenDelete)

		apiV1.GET("screen/studies/:studiesId/:operation", controller.ScreenGetStudiesOperation)
		apiV1.GET("screen/studies/:studiesId/series/:seriesId/:operation", controller.ScreenGetSeriesOperation)
		apiV1.GET("screen/studies/:studiesId/series/:seriesId/instances/:instanceId/:operation", controller.ScreenGetInstanceOperation)

		apiV1.POST("screen/studies/:studiesId/:operation", controller.ScreenPostStudiesOperation)
		apiV1.POST("screen/studies/:studiesId/series/:seriesId/:operation", controller.ScreenPostSeriesOperation)
		apiV1.POST("screen/studies/:studiesId/series/:seriesId/instances/:instanceId/:operation", controller.ScreenPostInstanceOperation)

		apiV1.GET("algo", controller.AlgoGet)
		apiV1.POST("algo", controller.AlgoPost)

		apiV1.GET("group", controller.GroupGet)

		apiV1.GET("raw", controller.RawDataGet)

		apiV1.GET("blockchain/nodelist", controller.GetBlockchainNodelist)
		apiV1.GET("blockchain/tps", controller.GetBlockchainTPS)
		apiV1.GET("blockchain/height", controller.GetBlockHeight)

		apiV1.POST("analysis/ct/:class/:mode", controller.AnalysisCtPost)

		aiApi := apiV1.Group("ai")
		aiApi.GET(":modal/:class/:algo/:aid", controller.AiAlgoGet)
		aiApi.GET(":modal/:class/:algo/:aid/:part", controller.AiAlgoGet)
		aiApi.POST(":modal/:class/:algo", controller.AiAlgoPost)
	}

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

	his := apiV1.Group("his")
	{
		his.Any("/:code", controller.ApiHisSearch)
	}

	strPort := fmt.Sprintf(":%d", config.SystemListenPort)

	if config.SystemUseHttps {
		log("i", "use https @ port", strPort)
		certPem := path.Join(config.PathApp, "certs/server.crt")
		certKey := path.Join(config.PathApp, "certs/server.key")
		if _, err := os.Stat(certKey); err != nil {
			panic(err.Error())
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

func TokenCheck(ctx *gin.Context) {
	key := ctx.Query("key")
	uid := -1
	valid := false

	if global.DebugGetFlag() && key == "ADMIN_BUAA" {
		log("w", "Debug ignore user check")
		ctx.Set("uid", 0)
		ctx.Next()
	} else if valid, uid = controller.CookieValidUid(ctx); !valid {
		//ctx.JSON(http.StatusNotAcceptable, "user invalid")
		ctx.Redirect(http.StatusFound, "/")
		ctx.Abort()
	} else {
		ctx.Set("uid", uid)
		ctx.Next()
	}
}
