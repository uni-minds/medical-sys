package module_router

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/controller"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"path"
	"time"
)

var Router *gin.Engine
var Instance *http.Server
var log *logger.Logger

func init() {
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".m3u8", "application/vnd.apple.mpegurl")
	// mime.AddExtensionType(".m3u8", "application/x-mpegurl")
	mime.AddExtensionType(".ts", "video/mp2t")
	// prevent on Windows with Dreamware installed, modified registry .css -> application/x-css
	// see https://stackoverflow.com/questions/22839278/python-built-in-server-not-loading-css
	mime.AddExtensionType(".css", "text/css; charset=utf-8")

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
}

func CreateServer(port int) *http.Server {
	Init()
	Instance = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           Router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return GetServer()
}

func GetServer() *http.Server {
	return Instance
}

func Init() {
	log = logger.NewLogger("ROUT")

	gin.DefaultWriter = logger.GetOutput()

	Router = gin.New()
	Router.Use(gin.Recovery(), gin.ErrorLogger())
	Router.Use(Cors())
	Router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		switch params.StatusCode {
		case 200:
			log.Log("t", fmt.Sprintf("%-4s 200 %s", params.Method, params.Path))
		default:
			log.Log("w", fmt.Sprintf("%-4s %d %s", params.Method, params.StatusCode, params.Path))
		}
		return ""
	}))

	paths := global.GetPaths()
	Router.Static("/plugins", path.Join(paths.Application, "web/plugins"))
	Router.Static("/application", paths.Application)
	Router.Static("/dist", path.Join(paths.Application, "web/dist"))
	Router.Static("/webapp", path.Join(paths.Application, "web/webapp"))
	Router.StaticFile("/favicon.ico", path.Join(paths.Application, "web/favicon.ico"))
	Router.StaticFile("/favicon-16x16.png", path.Join(paths.Application, "web/favicon-16x16.png"))
	Router.StaticFile("/favicon-32x32.png", path.Join(paths.Application, "web/favicon-32x32.png"))
	Router.StaticFile("/site.webmanifest", path.Join(paths.Application, "web/site.webmanifest"))
	Router.LoadHTMLGlob(path.Join(paths.Application, "web/templates/*"))

	Router.GET("/", controller.RootGet)
	// login logout forget register
	Router.GET("/:op", controller.RootGet)

	rUi := Router.Group("/ui", TokenCheck)
	{
		rUi.GET("/home", controller.UiHomeGet)
		rUi.GET("/manage/:class", controller.UiManageGetHandler)
		rUi.GET("/screen/studies/:studiesId/series/:seriesId/:operation", controller.UiScreenSeriesGet)
		rUi.GET("/labelsys/:mediaClass/:mediaIndex/:usertype", controller.UiLabeltoolGet)
		rUi.GET("/analysis", controller.UiAnalysisGet)
	}

	rMobi := Router.Group("/mobi", TokenCheck)
	{
		rMobi.GET("/", controller.MobiRoot)
		rMobi.GET("/device", controller.MobiGetDevice)
		rMobi.GET("/result/:pipeline", controller.MobiGetResult)
		rMobi.POST("/exec", controller.MobiMyExec)
	}

	apiV1 := Router.Group("/api/v1", TokenCheck)
	RouterApiVersion1(apiV1)

	Router.POST("/api/user/:op", controller.ApiUserPost)
}

func RouterApiVersion1(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	// /api/v1/labelsys/stream

	apiGroup.GET("/labelsys/stream/index/:index/:class/:op", controller.LabelsysGetStream)
	apiGroup.POST("/labelsys/stream/index/:index/:class/:op", controller.LabelsysPostStream)
	apiGroup.DELETE("/labelsys/stream/index/:index/:class/:op", controller.LabelsysDeleteStream)

	// /api/v1/medialist/
	apiGroup.GET("/medialist/:groupType/:groupId/:op", controller.MedialistGet)

	// /api/v1/user
	apiGroup.GET("user", controller.UserGet) // READ USER

	// /api/v1/sync
	apiGroup.POST("sync", controller.SyncPost)

	// /api/v1/media/:index/thumb
	apiGroup.GET("/media/index/:mediaIndex/:mediaOperate", controller.MediaGetOperation)

	// api pacs operation
	// studies
	apiGroup.POST("studies/:studiesId/:operation", controller.PostStudiesOperation)
	// series
	apiGroup.GET("studies/:studiesId/series/:seriesId/:operation", controller.SeriesGetOperation)
	apiGroup.POST("studies/:studiesId/series/:seriesId/:operation", controller.SeriesPostOperation)
	apiGroup.DELETE("studies/:studiesId/series/:seriesId/:operation", controller.SeriesDelOperation)
	// instance
	apiGroup.GET("studies/:studiesId/series/:seriesId/instances/:instanceId/:operation", controller.InstanceGetOperation)

	// /api/v1/algo
	apiGroup.GET("algo", controller.AlgoGet)
	apiGroup.POST("algo", controller.AlgoPost)

	// /api/v1/group
	apiGroup.GET("group", controller.GroupGet)

	// /api/v1/raw
	apiGroup.GET("raw", controller.RawDataGet)

	apiGroup.GET("blockchain/nodelist", controller.GetBlockchainNodelist)
	apiGroup.GET("blockchain/tps", controller.GetBlockchainTPS)
	apiGroup.GET("blockchain/height", controller.GetBlockHeight)

	apiGroup.POST("analysis/ct/:class/:mode", controller.AnalysisCtPost)

	// /api/v1/ai
	aiApi := apiGroup.Group("ai")
	{
		aiApi.GET(":modal/:class/:algo/:aid", controller.AiAlgoGet)
		aiApi.GET(":modal/:class/:algo/:aid/:part", controller.AiAlgoGet)
		aiApi.POST(":modal/:class/:algo", controller.AiAlgoPost)
	}

	// /api/v1/pacs
	pacs := apiGroup.Group("pacs")
	{
		pacs.Any("/", controller.PacsGet)
		pacs.Any("/:db", controller.PacsGetNodelist)
		pacs.Any("/:db/:node", controller.PacsSearchProxy)
		pacs.Any("/:db/:node/:studyuid", controller.PacsSearchProxy)
		pacs.Any("/:db/:node/:studyuid/:seriesuid", controller.PacsSearchProxy)
		pacs.Any("/:db/:node/:studyuid/:seriesuid/:objectuid", controller.PacsSearchProxy)
		//pacs.Any("/:db/:node/wado", controller.PacsWado)
	}

	// /api/v1/his
	his := apiGroup.Group("his")
	{
		his.Any("/:code", controller.ApiHisSearch)
	}

	return apiGroup
}

func TokenCheck(ctx *gin.Context) {
	//log("d", "token check")
	//key := ctx.Query("key")
	uid := -1
	valid := false
	valid, uid = controller.CookieValidUid(ctx)
	if valid {
		ctx.Set("uid", uid)
		ctx.Next()

	} else if global.FlagGetDebug() {
		log.Warn("debug mode: ignore user check")
		ctx.Set("uid", 1)
		ctx.Next()

	} else {
		log.Error(fmt.Sprintf("uid=%d: password check failed", uid))
		ctx.JSON(http.StatusNotAcceptable, controller.FailReturn(403, "user or password invalid"))
		ctx.Abort()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		//log("d", "cors origin:", origin)
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "172800")
		}

		if method == "OPTIONS" {
			//c.AbortWithStatus(http.StatusNoContent)
			c.JSON(http.StatusOK, "ok")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Error(fmt.Sprintf("%v", err))
			}
		}()
		c.Next()
	}
}
