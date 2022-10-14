package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/module"
)

func UIHomeGet(ctx *gin.Context) {
	_, uid := CookieValidUid(ctx)

	bgContent := fmt.Sprintf("您好，%s。\n"+
		"单击左侧菜单栏选择相应的功能。", module.UserGetRealname(uid))

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "首页 ｜ Medi-Sys",
		"page_id":    "index",
		"bg_content": bgContent,
	})
}
