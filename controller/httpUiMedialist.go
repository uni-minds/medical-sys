package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UIMedialistGet(ctx *gin.Context) {
	tp := ctx.Query("type")
	switch tp {
	case "us":
		ctx.HTML(http.StatusOK, "medialist_us.html", gin.H{
			"title":          "超声影像检索 | Medi-sys",
			"page_id":        "us-medialist",
			"custom_scripts": "/webapp/us/js/medialist-v2.js",
		})
		break

	case "ct":
		ctx.HTML(http.StatusOK, "medialist_ct.html", gin.H{
			"title":          "CT影像检索 | Medi-sys",
			"page_id":        "ct-medialist",
			"custom_scripts": "/webapp/ct/js/medialist.js",
		})
		break
	}
	return
}