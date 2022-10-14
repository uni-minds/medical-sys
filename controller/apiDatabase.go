package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"uni-minds.com/liuxy/medical-sys/tools"
)

type UriParams struct {
	Includefield string `form:"includefield"`
	StudyData    string `form:"StudyDate"`
	Offset       int    `form:"offset"`
	Limit        int    `form:"limit"`
}

func GetDatabaseDicomCtRsGroup(ctx *gin.Context) {
	class := ctx.Param("class")
	rsURI := strings.Split(ctx.Request.RequestURI, "/rs/")[1]

	var uriParams UriParams
	err := ctx.Bind(&uriParams)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn("检索异常"))
		return
	} else if (uriParams.Offset + uriParams.Limit) > 21 {
		ctx.JSON(http.StatusOK, FailReturn("演示系统禁止"))
		return
	}

	switch class {
	case "CTA":
		url := fmt.Sprintf("http://172.16.1.121:8080/dcm4chee-arc/aets/AS_RECEIVED/rs/%s", rsURI)
		_, t, bs, _ := tools.HttpGet(url)
		if bs == nil && t != "" {
			ctx.JSON(http.StatusOK, FailReturn(t))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(string(bs)))
		}
	case "CCTA":
		url := fmt.Sprintf("http://172.16.1.131:8080/dcm4chee-arc/aets/AS_RECEIVED/rs/%s", rsURI)
		_, t, bs, _ := tools.HttpGet(url)
		if bs == nil && t != "" {
			ctx.JSON(http.StatusOK, FailReturn(t))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(string(bs)))
		}

	default:
		ctx.JSON(http.StatusOK, FailReturn("Unknown type"))

	}
}

func GetDatabaseDicomCtWadoGroup(ctx *gin.Context) {
	var url string
	pURI := strings.Split(ctx.Request.RequestURI, "/wado/")[1]
	switch ctx.Param("class") {
	case "CTA":
		url = fmt.Sprintf("http://172.16.1.121:8080/dcm4chee-arc/aets/AS_RECEIVED/wado%s", pURI)
		fmt.Println(url)

	case "CCTA":
		url = fmt.Sprintf("http://172.16.1.131:8080/dcm4chee-arc/aets/AS_RECEIVED/wado%s", pURI)

	default:
		return
	}

	_, _, bs, _ := tools.HttpGet(url)
	ctx.Writer.Write(bs)
}
