package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UIManageGetHandler(ctx *gin.Context) {
	valid, _ := CookieValidUid(ctx)
	if !valid {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	action := ctx.Param("class")
	switch action {
	case "user":
		ctx.HTML(http.StatusOK, "manage-user.html", gin.H{
			"title":   "用户管理| Medi-sys",
			"page_id": "manage-user",
		})
		break

	case "group":
		ctx.HTML(http.StatusOK, "manage-group.html", gin.H{
			"page_id": "manage-group",
			"title":   "群组管理 ｜ Medi-sys",
		})
		break

	case "media":
		ctx.HTML(http.StatusOK, "manage-media.html", gin.H{
			"page_id": "manage-media",
			"title":   "Medi-sys | 媒体管理",
		})
		break

	case "blockchain":
		break

	case "upload":
		break
	}
	return
}
