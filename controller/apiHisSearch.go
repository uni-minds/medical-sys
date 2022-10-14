package controller

import (
	"gitee.com/uni-minds/medical-sys/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiHisSearch(ctx *gin.Context) {
	// /api/v1/his/:code
	code := ctx.Param("code")
	data, err := database.BridgeGetHisDatabaseRetrieve(code)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(data))
	}
}
