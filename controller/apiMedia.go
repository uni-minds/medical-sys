package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/manager"
	"uni-minds.com/liuxy/medical-sys/module"
)

type mediaInfoForJsGrid struct {
	Mid       int                      `json:"mid"`
	MediaHash string                   `json:"media"`
	Name      string                   `json:"name"`
	View      string                   `json:"view"`
	Duration  float64                  `json:"duration"`
	Frames    int                      `json:"frames"`
	Authors   labelInfoForJsGridButton `json:"authors"`
	Reviews   labelInfoForJsGridButton `json:"reviews"`
	Memo      string                   `json:"memo"`
}

type labelInfoForJsGridButton struct {
	Realname string `json:"realname"`
	Tips     string `json:"tips"`
	Status   string `json:"status"`
}

type medialistForJsGrid struct {
	Data       []mediaInfoForJsGrid `json:"data"`
	ItemsCount int                  `json:"itemsCount"`
}

func MediaGetHandler(ctx *gin.Context) {
	valid, uid := CookieValidUid(ctx)
	if !valid {
		return
	}

	action := ctx.Query("action")
	switch action {
	case "getlist":
		gid, err := strconv.Atoi(ctx.Query("gid"))
		if err != nil {
			return
		}

		var callback medialistForJsGrid

		page, _ := strconv.Atoi(ctx.Query("page"))
		if page <= 0 {
			page = 1
		}
		count, _ := strconv.Atoi(ctx.Query("count"))
		if count <= 0 {
			count = 20
		}

		var mids []int
		order := ctx.Query("order")
		field := ctx.Query("field")
		if field != "" {
			mids = module.UserGetGroupMediaSelector(uid, gid, field, order)
			callback.ItemsCount = len(mids)
		} else {
			mids = module.UserGetGroupMedia(uid, gid)
			callback.ItemsCount = len(mids)
		}

		index := (page - 1) * count

		mdata := make([]mediaInfoForJsGrid, 0)
		if callback.ItemsCount > index {
			_ = database.UserSetLastStatus(uid, gid, page)
			mids = mids[index:]
			if len(mids) >= count {
				mids = mids[0:count]
			}

			for _, mid := range mids {
				mediaSummary, err := module.MediaGetSummary(mid)
				if err != nil {
					log.Println("E get mediaSummary", mid, err.Error())
					continue
				}

				labelSummary, err := module.LabelGetSummary(mediaSummary.Hash)
				var authorData, reviewData labelInfoForJsGridButton
				if err == nil {
					authorData = labelInfoForJsGridButton{
						Realname: labelSummary.AuthorRealname,
						Tips:     labelSummary.AuthorTips,
						Status:   labelSummary.AuthorProgress,
					}

					reviewData = labelInfoForJsGridButton{
						Realname: labelSummary.ReviewRealname,
						Tips:     labelSummary.ReviewTips,
						Status:   labelSummary.ReviewProgress,
					}
				}

				mdata = append(mdata, mediaInfoForJsGrid{
					Mid:       mid,
					MediaHash: mediaSummary.Hash,
					Name:      mediaSummary.DisplayName,
					View:      mediaSummary.Views,
					Duration:  mediaSummary.Duration,
					Frames:    mediaSummary.Frames,
					Authors:   authorData,
					Reviews:   reviewData,
					Memo:      mediaSummary.Memo,
				})
			}
		}
		callback.Data = mdata
		//jc, _ := json.Marshal(callback)
		//ctx.JSON(http.StatusOK, SuccessReturn(string(jc)))
		ctx.JSON(http.StatusOK, SuccessReturn(callback))
		return

	case "play":
		mediaHash := ctx.Query("media")
		if len(mediaHash) < 32 {
			return
		}
		fp := module.MediaGetRealpath(mediaHash, uid)
		ctx.File(fp)
		return

	case "getlock":
		mediaHash := ctx.Query("media")

		if mediaHash == "ALL" {
			ctx.JSON(http.StatusOK, SuccessReturn(manager.MediaAccessLockList()))
			return
		}

		if len(mediaHash) < 32 {
			return
		}

		status, err := manager.MediaAccessGetLock(mediaHash)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(status))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(status))
		}

	case "setlock":
		mediaHash := ctx.Query("media")
		if len(mediaHash) < 32 {
			return
		}

		tp := ctx.Query("type")
		switch tp {
		case "author", "review":
			status, err := manager.MediaAccessSetLock(mediaHash, uid, tp)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(status))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(status))
			}
		}

	case "unlock":
		mediaHash := ctx.Query("media")
		if len(mediaHash) < 32 {
			return
		}

		manager.MediaAccessUnlock(mediaHash, uid, true)
		ctx.JSON(http.StatusOK, SuccessReturn("OK"))

	default:
		ctx.JSON(http.StatusOK, FailReturn(action))

	}
}
