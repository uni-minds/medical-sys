package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UIMedialistGetHandler(ctx *gin.Context) {
	valid, _ := CookieValidUid(ctx)
	if !valid {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	tp := ctx.Query("type")
	switch tp {
	case "us":
		ctx.HTML(http.StatusOK, "medialist.html", gin.H{
			"title":          "超声影像检索 | Medi-sys",
			"page_id":        "us-medialist",
			"custom_scripts": "/webapp/ultrasonic/js/ultrasonic_medialist-v2.js",
		})
		break

	case "ct":
		ctx.HTML(http.StatusOK, "medialist_ct.html", gin.H{
			"page_id": "manage-group",
			"title":   "CT影像管理 | Medi-sys",
		})
		break
	}
	return
}
