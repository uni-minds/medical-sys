package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/manager"
	"uni-minds.com/liuxy/medical-sys/module"
	"uni-minds.com/liuxy/medical-sys/tools"
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

func MediaGet(ctx *gin.Context) {
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
					log("e", "E get mediaSummary", mid, err.Error())
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
		//ctx.JSON(http.StatusOK, tools.SuccessReturn(string(jc)))
		ctx.JSON(http.StatusOK, SuccessReturn(callback))
		return

	case "play":
		mediaHash := ctx.Query("media")
		if len(mediaHash) < 32 {
			return
		}
		fogv := module.MediaGetRealpath(mediaHash, uid)
		if _, err := os.Stat(fogv); err != nil {
			fmt.Println("E;media file not found!", fogv)
			return
		}

		t := ctx.Query("type")
		switch t {
		case "mp4":
			fmp4 := strings.Replace(fogv, ".ogv", ".mp4", 1)
			if _, err := os.Stat(fmp4); err != nil {
				fmp4 = fmt.Sprintf("%s.mp4", fogv)
				if _, err := os.Stat(fmp4); err != nil {
					fmt.Printf("ffmpeg convert: %s => %s\n", fogv, fmp4)
					if err := tools.FFmpegToH264(fogv, fmp4); err != nil {
						fmt.Println("E;ffmpeg convert:", err.Error())
						return
					}
					fmt.Println("ffmpeg convert finish.")
				}
			}

			ctx.File(fmp4)
			return

		default:
			ctx.File(fogv)

		}

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
			ctx.JSON(http.StatusOK, FailReturn(400, status))
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
				ctx.JSON(http.StatusOK, FailReturn(400, status))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(status))
			}
		}

	case "setunlock":
		mediaHash := ctx.Query("media")
		if len(mediaHash) < 32 {
			return
		}

		manager.MediaAccessUnlock(mediaHash, uid, true)
		ctx.JSON(http.StatusOK, SuccessReturn("OK"))

	default:
		ctx.JSON(http.StatusOK, FailReturn(400, action))

	}
}
