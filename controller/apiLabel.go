package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"uni-minds.com/liuxy/medical-sys/global"
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

func LabelGet(ctx *gin.Context) {
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

	case "author", "review":
		switch ctx.Query("selector") {
		case "memo":
			memo := module.UserGetMediaMemo(uid, module.UserGetMid(uid, mediaHash))
			ctx.JSON(http.StatusOK, SuccessReturn(memo))
			return

		case "full":
			ld := module.LabelGetJson(mediaHash)
			ctx.JSON(http.StatusOK, SuccessReturn(ld))
			return
		}
	}
}

func LabelPost(ctx *gin.Context) {
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
			} else if ldata.MediaHash != mediaHash {
				ctx.JSON(http.StatusOK, FailReturn("上传数据特征异常"))
			} else if err = module.UserSetMediaMemo(uid, mid, ldata.Data); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn("OK"))
			}
			return

		case "reject":
			err := module.LabelSubmitReview(mediaHash, uid, "reject")
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else {
				log.Println("User label confirmed.")
				ctx.JSON(http.StatusOK, SuccessReturn("exit"))
			}
			return

		case "confirm":
			log.Println("User label confirm", uid, mid)
			err := module.LabelSubmitReview(mediaHash, uid, "confirm")
			if err != nil {
				log.Println("E User label confirm write", err.Error())
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else {
				log.Println("User label confirmed.")
				ctx.JSON(http.StatusOK, SuccessReturn("exit"))
			}
			return

		case "full":
			var ldata LabelData
			if err := ctx.BindJSON(&ldata); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))

			} else if err = module.LabelUpdateReview(ldata.Data, mediaHash, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))

			} else {
				fmt.Println("Reviewer json import:", ldata)
				ctx.JSON(http.StatusOK, SuccessReturn(fmt.Sprintf("同步成功 @ %s", time.Now().Format(global.TimeFormat))))

			}
			return

		default:
			ctx.JSON(http.StatusOK, FailReturn("功能不可用"))
			return
		}

	case "author":
		switch ctx.Query("selector") {
		case "memo":
			var ldata LabelData
			if err := ctx.BindJSON(&ldata); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else if err = module.UserSetMediaMemo(uid, mid, ldata.Data); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(fmt.Sprintf("同步成功 @ %s", time.Now().Format(global.TimeFormat))))
			}
			return

		case "full":
			var ldata LabelData
			if err := ctx.BindJSON(&ldata); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else if err := module.LabelUpdateAuthor(ldata.Data, mediaHash, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else {
				fmt.Println("Author json import:")
				fmt.Println(ldata.Data)
				fmt.Println("---")
				ctx.JSON(http.StatusOK, SuccessReturn(fmt.Sprintf("同步成功 @ %s", time.Now().Format(global.TimeFormat))))
			}
			return

		case "submit":
			// 提交
			if err := module.LabelSubmitAuthor(mediaHash, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn("exit"))
			}
			return

		default:
			ctx.JSON(http.StatusOK, FailReturn("功能不可用"))
			return
		}
	}
}

func LabelDel(ctx *gin.Context) {

}
