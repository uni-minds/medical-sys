/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiMedia.go
 */

package controller

import (
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/manager"
	"gitee.com/uni-minds/medical-sys/module"
	"gitee.com/uni-minds/medical-sys/tools"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const DICOM_TYPE_US_ID = "1.2.276.0.26.1.1.1.2"

type mediaInfoForJsGrid struct {
	Mid       string                   `json:"mid"`
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
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	switch ctx.Query("action") {
	case "getlist":
		var callback medialistForJsGrid

		gid, err := strconv.Atoi(ctx.Query("gid"))
		if err != nil {
			return
		}
		page, _ := strconv.Atoi(ctx.Query("page"))
		if page <= 0 {
			page = 1
		}
		count, _ := strconv.Atoi(ctx.Query("count"))
		if count <= 0 {
			count = 20
		}
		order := ctx.Query("order")
		field := ctx.Query("field")

		index := (page - 1) * count

		groupType := ctx.Query("type")
		switch groupType {
		case "mid", "label_media":
			var mediaIndex []string
			if field != "" {
				mediaIndex = module.UserGetGroupContainsSelector(uid, gid, field, order)
				callback.ItemsCount = len(mediaIndex)
			} else {
				mediaIndex, _, _ = module.UserGetGroupContains(uid, gid)
				callback.ItemsCount = len(mediaIndex)
			}

			mdata := make([]mediaInfoForJsGrid, 0)
			if callback.ItemsCount > index {
				_ = database.UserSetLastStatus(uid, gid, page)
				mediaIndex = mediaIndex[index:]
				if len(mediaIndex) >= count {
					mediaIndex = mediaIndex[0:count]
				}

				for _, mid := range mediaIndex {
					id, _ := strconv.Atoi(mid)
					mediaSummary, err := module.MediaGetSummary(id)
					if err != nil {
						log("e", "E get mediaSummary", mid, err.Error())
						continue
					}

					labelSummary, _, _, err := module.LabelGetSummary(mediaSummary.Hash)
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

		case "label_dicom", "instanceid":
			server := database.BridgeGetPacsServerHandler()
			instanceIds, _, err := module.UserGetGroupContains(uid, gid)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(403, err.Error()))
				return
			}

			callback.ItemsCount = len(instanceIds)
			mdata := make([]mediaInfoForJsGrid, 0)

			if callback.ItemsCount > index {
				_ = database.UserSetLastStatus(uid, gid, page)
				end := index + count
				if end > len(instanceIds) {
					end = len(instanceIds)
				}

				for _, instanceId := range instanceIds[index:end] {
					info, err := server.FindInstanceByIdLocal(instanceId)
					if err != nil {
						log("e", err.Error())
						continue
					}

					var authorData, reviewData labelInfoForJsGridButton
					if labelSummary, _, _, err := module.LabelGetSummary(instanceId); err == nil {
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
						Mid:       instanceId,
						MediaHash: instanceId,
						Name:      instanceId,
						View:      info.LabelView,
						Duration:  info.Duration,
						Frames:    info.Frames,
						Authors:   authorData,
						Reviews:   reviewData,
						Memo:      module.LabelGetMemo(instanceId),
					})
				}
			}
			callback.Data = mdata
		}

		ctx.JSON(http.StatusOK, SuccessReturn(callback))
		return
	}
}

func MediaPreOperation(ctx *gin.Context) {
	mediaIndex := ctx.Param("mediaIndex")
	if mediaIndex == "" {
		ctx.JSON(http.StatusOK, FailReturn(400, "empty index"))
		ctx.Abort()
	}
	ctx.Next()
}

func MediaGetOperation(ctx *gin.Context) {
	//uid := -1
	//if value, exists := ctx.Get("uid"); !exists {
	//	return
	//} else {
	//	uid = value.(int)
	//}

	mediaIndex := ctx.Param("mediaIndex")
	op := ctx.Param("op")

	switch op {
	case "check":
		tp := ctx.Query("type")
		switch tp {
		case "ogv":
			ctx.JSON(http.StatusOK, SuccessReturn("OK"))
		case "mp4":
		case "jpg":
		case "png":
		default:
			ctx.JSON(http.StatusOK, FailReturn(403, fmt.Sprintf("unknown type: %s", tp)))

		}

	case "video":
		tp := ctx.Query("type")
		if tp == "" {
			tp = "video"
		}

		if bs, _, err := module.PacsGetInstanceMedia(mediaIndex, tp); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(403, err.Error()))
		} else {
			ctx.Writer.Write(bs)
		}

	case "image":
		bs, _, err := module.PacsGetInstanceMedia(mediaIndex, "image")

		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
		} else {
			ctx.Writer.Write(bs)
		}

	case "lock":
		switch mediaIndex {
		case "*":
			ctx.JSON(http.StatusOK, SuccessReturn(manager.MediaAccessLockList()))
		default:
			status, err := manager.MediaAccessGetLock(mediaIndex)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, status))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(status))
			}
		}

	case "thumb":
		switch ctx.Query("size") {
		default:
			bs, err := module.PacsGetInstanceThumb(mediaIndex)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
			} else {
				ctx.Writer.Write(bs)
			}
		}
	}
}

func MediaPostOperation(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	mediaIndex := ctx.Param("mediaIndex")
	op := ctx.Param("op")

	switch op {
	case "lock":
		status, err := manager.MediaAccessSetLock(mediaIndex, uid, "")
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(status))
		}
	}
}

func MediaDeleteOperation(ctx *gin.Context) {
	uid := -1
	if value, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = value.(int)
	}

	mediaIndex := ctx.Param("mediaIndex")
	op := ctx.Param("op")

	switch op {
	case "lock":
		manager.MediaAccessUnlock(mediaIndex, uid, true)
		ctx.JSON(http.StatusOK, SuccessReturn("OK"))
	}
}

func MediaGetData(mediaIndex string, uid int, mediaType string) ([]byte, error) {
	if strings.Contains(mediaIndex, DICOM_TYPE_US_ID) {
		return module.InstanceGetVideo(mediaIndex, uid)

	} else if len(mediaIndex) < 32 {
		return nil, errors.New("unknown mediaIndex")
	} else {
		fogv := module.MediaGetRealpath(mediaIndex, uid)
		if _, err := os.Stat(fogv); err != nil {
			return nil, err
		}

		var filepath string
		switch mediaType {
		case "mp4":
			fmp4 := strings.Replace(fogv, ".ogv", ".mp4", 1)
			if _, err := os.Stat(fmp4); err != nil {
				fmp4 = fmt.Sprintf("%s.mp4", fogv)
				if _, err := os.Stat(fmp4); err != nil {
					fmt.Printf("ffmpeg convert: %s => %s\n", fogv, fmp4)
					if err := tools.FfConv(fogv, fmp4, "h264"); err != nil {
						fmt.Println("E;ffmpeg convert:", err.Error())
						return nil, err
					}
					fmt.Println("ffmpeg convert finish.")
				}
			}

			filepath = fmp4

		default:
			filepath = fogv
		}

		if fp, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm); err != nil {
			return nil, err
		} else if bs, err := ioutil.ReadAll(fp); err != nil {
			fp.Close()
			return nil, err
		} else {
			fp.Close()
			return bs, nil
		}
	}
}
