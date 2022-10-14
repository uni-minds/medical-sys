/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: apiLabelsysGetStream.go
 * Date: 2022/09/15 13:05:15
 */

// ECode = 40200x

package controller

import (
	"fmt"
	pacs_dcm4chee "gitee.com/uni-minds/bridge-pacs/dcm4chee"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/manager"
	"gitee.com/uni-minds/medical-sys/module"
	"gitee.com/uni-minds/utils/media"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func LabelsysGetStream(ctx *gin.Context) {
	// apiGroup.GET("/labelsys/stream/index/:index/:class/:op", controller.LabelsysGetStream)
	mediaUUID := ctx.Param("index")
	opClass := ctx.Param("class")
	opType := ctx.Param("op")

	info, err := module.MediaGet(mediaUUID)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(402000, err.Error()))
		return
	}

	switch opClass {
	case "info":
		log("d", "labelsys info", mediaUUID, opClass, opType)
		//mediaUrl := fmt.Sprintf("/storage/%s/%s/media", mediaClass, mediaUUID)
		switch opType {
		// /api/v1/labelsys/stream/index/INDEX/info/full
		case "full":
			var mi MediaInfo
			switch info.MediaType {
			case pacs_dcm4chee.MEDIA_TYPE_IMAGE:
				mi = MediaInfo{
					MediaDuration: 0,
					MediaFrames:   0,
					MediaFps:      0,
					MediaHeight:   info.Height,
					MediaWidth:    info.Width,
					MediaIndex:    mediaUUID,
					MediaURL:      fmt.Sprintf("/api/v1/media/index/%s/ogv", mediaUUID),
					MustHLS:       false,
				}

			case pacs_dcm4chee.MEDIA_TYPE_MULTI_FRAME:
				// instance id video
				mi = MediaInfo{
					MediaDuration: info.Duration,
					MediaFrames:   info.Frames,
					MediaFps:      info.Fps,
					MediaHeight:   info.Height,
					MediaWidth:    info.Width,
					MediaIndex:    mediaUUID,
					MediaURL:      fmt.Sprintf("/api/v1/media/index/%s/mp4", mediaUUID),
					MustHLS:       false,
				}

			default:
				mediaFrame := info.Frames
				mediaDuration := info.Duration
				mediaFps := info.Fps

				if mediaFrame == 0 || mediaDuration == 0 || mediaFps == 0 {
					if mediaDur, err := media.GetDuration(info.Path); err != nil {
						ctx.JSON(http.StatusOK, FailReturn(402001, err.Error()))
					} else {
						mediaDuration = mediaDur.Seconds()
					}
					mediaFps, _ = media.GetFps(info.Path)
					mediaFrame = int(mediaDuration * mediaFps)
				}

				var targetUrl string
				var useHls bool

				switch path.Ext(info.Path) {
				case ".ogv", ".mp4":
					log("d", "force use mp4")
					targetUrl = fmt.Sprintf("/api/v1/media/index/%s/mp4", mediaUUID)
					useHls = false

				case ".m3u8":
					targetUrl = fmt.Sprintf("/api/v1/media/index/%s/video.m3u8", mediaUUID)
					useHls = true

				}

				mi = MediaInfo{
					MediaDuration: mediaDuration,
					MediaFrames:   mediaFrame,
					MediaFps:      mediaFps,
					MediaHeight:   info.Height,
					MediaWidth:    info.Width,
					MediaIndex:    mediaUUID,
					MediaURL:      targetUrl,
					MustHLS:       useHls,
				}
			}

			ctx.JSON(http.StatusOK, SuccessReturn(mi))
		default:
			ctx.JSON(http.StatusOK, FailReturn(402002, "unknown op: "+opType))
		}
		return

	case "label":
		switch opType {
		case "data", "review", "author":
			labelUUID, _, _, _ := module.MediaGetFirstLabeler(info.MediaUUID)

			if labelUUID != "" {
				labelstr := module.LabelGetJson(labelUUID)
				ctx.JSON(http.StatusOK, SuccessReturn(labelstr))
			} else {
				ctx.JSON(http.StatusOK, FailReturn(10001, "no label"))
			}

		case "authors":
			if len(info.LabelAuthors) > 0 {
				ctx.JSON(http.StatusOK, SuccessReturn(info.LabelAuthors))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(nil))
			}

		case "reviewers":
			if len(info.LabelReviewers) > 0 {
				ctx.JSON(http.StatusOK, SuccessReturn(info.LabelReviewers))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(nil))
			}

		case "crf":
			view := info.CrfDefine
			if view != "" {
				ctx.JSON(http.StatusOK, SuccessReturn(global.LoadCrfViewData(view)))
			} else {
				ctx.JSON(http.StatusOK, FailReturn(402004, "no crf define"))
			}
			return

		case "memo":
			ctx.JSON(http.StatusOK, SuccessReturn(info.Memo))

		case "summary":
			uuids, err := module.MediaGetLabelUUIDs(info.MediaUUID)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(30002, err.Error()))
				return
			}

			if len(uuids) > 0 {
				summary, _, _, err := module.LabelGetSummary(uuids[0])
				if err != nil {
					ctx.JSON(http.StatusOK, FailReturn(30001, err.Error()))
				} else {
					ctx.JSON(http.StatusOK, SuccessReturn(summary))
				}
			} else {
				// 无标注
				ctx.JSON(http.StatusOK, FailReturn(30001, "no author"))
			}

		default:
			if len(opType) != 32 {
				ctx.JSON(http.StatusOK, FailReturn(402005, "unknow label operation"))
			} else if data, err := module.LabelGet(opType); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(402005, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn(data))
			}
		}

		return

		// lock/data GET|PUT|DEL

	case "lock":
		switch opType {
		case "data":
			switch mediaUUID {
			case "*":
				ctx.JSON(http.StatusOK, SuccessReturn(manager.MediaAccessLockList()))
			default:
				status, err := manager.MediaAccessGetLock(mediaUUID)
				if err != nil {
					// 无锁
					ctx.JSON(http.StatusOK, FailReturn(400, "unlock"))
				} else {
					// 有锁
					ctx.JSON(http.StatusOK, SuccessReturn(status))
				}
			}
			return
		}

	case "screen":
		ctx.JSON(http.StatusOK, FailReturn(402007, "unknown screen operation"))

	default:
		ctx.JSON(http.StatusOK, FailReturn(402008, "unknown operate class"))

	}
}
