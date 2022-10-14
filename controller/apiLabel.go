package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"uni-minds.com/liuxy/medical-sys/module"
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

	mediaHash := ctx.Query("media")

	switch ctx.Query("action") {
	case "getrealname":
		authorName, reviewName := module.LabelGetRealname(mediaHash)
		ctx.JSON(http.StatusOK, SuccessReturn([]string{authorName, reviewName}))

	case "summary":
		var err error
		summary, err := module.LabelGetSummary(mediaHash)
		if err != nil {
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

		switch ctx.Query("selector") {
		case "memo":
			var ldata LabelData
			err := ctx.BindJSON(&ldata)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}
			if ldata.MediaHash != mediaHash {
				ctx.JSON(http.StatusOK, FailReturn("上传数据特征异常"))
				return
			}
			err = module.UserSetMediaMemo(uid, mid, ldata.Data)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}

		case "reject":
			err := module.LabelSubmitReview(mediaHash, uid, "reject")
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			}

		case "confirm":
			log.Println("User label confirm", uid, mid)
			err := module.LabelSubmitReview(mediaHash, uid, "confirm")
			if err != nil {
				log.Println("E User label confirm write", err.Error())
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else {
				log.Println("User label confirmed.")
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
			err := module.LabelUpdateAuthor(ldata.Data, mediaHash, uid)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}

		case "submit":
			// 提交
			err := module.LabelSubmitAuthor(mediaHash, uid)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
				return
			}

		}

		ctx.JSON(http.StatusOK, SuccessReturn("OK"))
	}
}
