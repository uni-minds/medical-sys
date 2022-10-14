package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"uni-minds.com/medical-sys/module"
)

type LabelData struct {
	MediaHash string `json:"media"`
	Data      string `json:"data"`
	Direction string `json:"direction"`
}

type LabelInfoForButton struct {
	TextBackG string `json:"textbackg"`
	TextHover string `json:"texthover"`
	Tips      string `json:"tips"`
}

func LabelGetHandler(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if !valid {
		ctx.JSON(http.StatusOK, FailReturn(ETokenInvalid))
		return
	}

	labelHash := ctx.Query("label")
	mediaHash := ctx.Query("media")

	switch ctx.Query("action") {
	case "getrealname":
		realname := module.LabelGetRealname(labelHash)
		ctx.JSON(http.StatusOK, SuccessReturn(realname))

	case "summary":
		var summary LabelInfoForButton
		var err error
		summary.TextBackG, summary.TextHover, summary.Tips, err = module.LabelGetSummary(labelHash, uid)
		if err != nil {
			log.Println("Label summary E:", err.Error())
			ctx.JSON(http.StatusOK, FailReturn(err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(summary))
		}

	case "review":
		switch ctx.Query("selector") {
		case "memo":
			memo := module.UserGetMediaMemo(uid, module.UserGetMid(uid, mediaHash))
			ctx.JSON(http.StatusOK, SuccessReturn(memo))
			return
		}

	case "author":
		switch ctx.Query("selector") {
		case "memo":
			memo := module.UserGetMediaMemo(uid, module.UserGetMid(uid, mediaHash))
			ctx.JSON(http.StatusOK, SuccessReturn(memo))
			return
		}
	}
}

func LabelPostHandler(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if !valid {
		ctx.JSON(http.StatusOK, FailReturn(ETokenInvalid))
		return
	}

	action := ctx.Query("action")
	mediaHash := ctx.Query("media")
	mid := module.UserGetMid(uid, mediaHash)

	switch action {
	case "review":
		var ldata LabelData
		err := ctx.BindJSON(&ldata)
		if err != nil || ldata.MediaHash != mediaHash {
			ctx.JSON(http.StatusOK, FailReturn(err.Error()))
		}

		switch ctx.Query("selector") {
		case "memo":
			err := module.UserSetMediaMemo(uid, mid, ldata.Data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}

		case "full":
			err := module.LabelSetReviewerJson(mid, uid, ldata.Data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}
		}
		ctx.JSON(http.StatusOK, SuccessReturn("OK"))

	case "author":
		var ldata LabelData
		err := ctx.BindJSON(&ldata)
		if err != nil || ldata.MediaHash != mediaHash {
			ctx.JSON(http.StatusOK, FailReturn(err.Error()))
		}

		switch ctx.Query("selector") {
		case "memo":
			err := module.UserSetMediaMemo(uid, mid, ldata.Data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}

		case "full":
			err := module.LabelSetAuthorJson(mid, uid, ldata.Data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}
		}
		ctx.JSON(http.StatusOK, SuccessReturn("OK"))
	}
}
