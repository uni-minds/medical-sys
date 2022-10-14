package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
	"time"
	"uni-minds.com/medical-sys/controller"
	"uni-minds.com/medical-sys/database"
	"uni-minds.com/medical-sys/global"
	"uni-minds.com/medical-sys/module"
	"uni-minds.com/medical-sys/upgrade"
)

const serverPort = ":8443"

const webRoot = "application/web"
const mediaRoot = "application/media"

func initVariables() {
	global.SetCookieMaxAge(24 * int(time.Hour.Seconds()))
	global.SetMediaRoot(mediaRoot)
	global.SetAppSettings(global.AppSettings{
		EnableUserRegister: true,
	})
	global.SetUserRegCode("bjkr2020")
}

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
	//initVariables()
	//upgrade.UpgradeParseMediaData()
	//prerun()
	//providerWeb(serverPort)
	convert()
}

func prerun() {
	//upgrade.UpgradeImportUsers()
	upgrade.UpgradeImportGroupMedia("common")
	module.GroupCreate("4ap_review_1", "四腔心审核#1", "4ap_review_1_liu")
	module.GroupCreate("4ap_review_2", "四腔心审核#2", "4ap_review_1_gu")
	module.GroupCreate("4ap_review_3", "四腔心审核#3", "4ap_review_1_zhang")
	module.GroupCreate("4ap_review_4", "四腔心审核#4", "4ap_review_1_han")
	module.GroupAddUser(module.GroupGetGid("4ap_review_1"), 3, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_2"), 4, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_3"), 5, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_4"), 6, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_1"), 25, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_2"), 25, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_3"), 25, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_4"), 25, "leader")
	module.GroupAddUser(module.GroupGetGid("common"), 3, "leader")
	module.GroupAddUser(module.GroupGetGid("common"), 4, "leader")
	module.GroupAddUser(module.GroupGetGid("common"), 5, "leader")
	module.GroupAddUser(module.GroupGetGid("common"), 6, "leader")
	module.UserSetPassword(3, "bjkr@2020")
	module.UserSetPassword(4, "bjkr@2020")
	module.UserSetPassword(5, "bjkr@2020")
	module.UserSetPassword(6, "bjkr@2020")

	module.GroupAddUser(module.GroupGetGid("4ap_review_1"), 1, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_2"), 1, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_3"), 1, "leader")
	module.GroupAddUser(module.GroupGetGid("4ap_review_4"), 1, "leader")
}

func grouping() {
	gids := []int{9, 10, 11, 12}
	gi := 0
	tstd, _ := time.Parse(global.TimeFormat, "2020-01-15 00:00:00")
	for mid := 1; mid < 326; mid++ {
		lis, err := database.LabelGetAll(mid, 0, global.LabelTypeAuthor)
		if err != nil {
			continue
		}

		flag := false

		for i := 0; i < len(lis); i++ {
			tcreat, _ := time.Parse(global.TimeFormat, lis[i].CreateTime)
			if tcreat.Sub(tstd) < 0 {
				flag = true
				break
			}
		}

		if flag {
			if gi == 4 {
				gi = 0
			}
			database.GroupAddMedia(gids[gi], mid)
			gi++
		}
	}
}

func convert() {
	lis, err := database.LabelGetAll(0, 0, global.LabelTypeAuthor)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, li := range lis {
		b := database.LabelsInfo{
			Lid:               i + 1,
			Progress:          2,
			AuthorUid:         li.Uid,
			ReviewUid:         0,
			MediaHash:         strconv.Itoa(li.Mid),
			Data:              li.DataBackup,
			Version:           1,
			Frames:            li.Frames,
			Counts:            li.Counts,
			TimeAuthorStart:   li.CreateTime,
			TimeAuthorSubmit:  li.ModifyTime,
			TimeReviewStart:   "",
			TimeReviewConfirm: "",
			Memo:              li.Memo,
		}
		database.LabelsCreate(b)
	}

	lis, err = database.LabelGetAll(0, 0, global.LabelTypeReview)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, li := range lis {
		b, err := database.LabelsGet(li.Mid)
		if err != nil {
			fmt.Println(err.Error())
		}
		b.Data = li.DataBackup
		b.ReviewUid = li.Uid
		b.TimeReviewStart = li.CreateTime
		b.TimeReviewSubmit = li.ModifyTime
		database.LabelsUpdate(b)
	}
}
