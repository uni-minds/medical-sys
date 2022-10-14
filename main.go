package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"time"
	"uni-minds.com/liuxy/medical-sys/controller"
	"uni-minds.com/liuxy/medical-sys/global"
)

const serverPort = ":8443"
const webRoot = "application/web"
const mediaRoot = "application/media"

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func providerWeb(port string) {
	router := gin.Default()

	//router.Use(Cors())
	router.Static("/build", path.Join(webRoot, "build"))
	router.Static("/plugins", path.Join(webRoot, "plugins"))
	router.Static("/pages", path.Join(webRoot, "pages"))
	router.Static("/dist", path.Join(webRoot, "dist"))
	router.Static("/webapp", path.Join(webRoot, "webapp"))
	router.StaticFile("/favicon.ico", path.Join(webRoot, "favicon.ico"))

	router.LoadHTMLGlob(path.Join(webRoot, "templates", "**/*"))

	router.GET("/", controller.UIRootGetHandler)

	//user interface
	ui := router.Group("/ui")
	{
		ui.GET("home", controller.UIHomeGetHandler)
		ui.GET("register", controller.UIRegisterGetHandler)
		ui.GET("manage/:class", controller.UIManageGetHandler)
		ui.GET("medialist", controller.UIMedialistGetHandler)
		ui.GET("labeltool", controller.UILabeltoolGetHandler)
		ui.GET("logout", controller.UILogoutGetHandler)
		ui.GET("import", controller.UIImportMedia)
	}

	//region API
	api := router.Group("api")
	{
		api.POST("register", controller.RegisterPost)  // register
		api.POST("login", controller.LoginPostHandler) // Login
		// DEBUG Only
		// https://localhost:8443/api/login?goldenkey=Uni-Ledger-RIS&user=$
		api.GET("login", controller.LoginGetHandler)
		// DEBUG Only

		api.POST("user", controller.UserPostHandler)     // CREATE USER
		api.GET("user", controller.UserGetHandler)       // READ USER
		api.DELETE("user", controller.UserDeleteHandler) // REMOVE USER
		api.PUT("user", controller.UserPutHandler)       // UPDATE USER

		api.POST("media")
		api.GET("media", controller.MediaGetHandler)
		api.DELETE("media")
		api.PUT("media")

		api.POST("label", controller.LabelPostHandler)
		api.GET("label", controller.LabelGetHandler)
		api.DELETE("label")
		api.PUT("label")

		api.POST("group")
		api.GET("group", controller.GroupGetHandler)
		api.DELETE("group")
		api.PUT("group")
	}
	//endregion

	router.RunTLS(port, "./application/cert.pem", "./application/cert.key")
}

func main() {
	global.SetCookieMaxAge(24 * int(time.Hour.Seconds()))
	global.SetMediaRoot(mediaRoot)
	global.SetAppSettings(global.AppSettings{
		EnableUserRegister: true,
	})
	global.SetUserRegCode("bjkr2020")
	providerWeb(serverPort)
}
