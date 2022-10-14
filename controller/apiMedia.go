/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiMedia.go
 */

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
	"os"
	"path"
	"strings"
)

type mediaInfoForJqGrid struct {
	Mid       string                     `json:"mid"`
	MediaUUID string                     `json:"media"`
	Name      string                     `json:"name"`
	View      string                     `json:"view"`
	Duration  float64                    `json:"duration"`
	Frames    int                        `json:"frames"`
	Authors   []labelInfoForJsGridButton `json:"authors"`
	Reviews   []labelInfoForJsGridButton `json:"reviews"`
	Memo      string                     `json:"memo"`
}

type labelInfoForJsGridButton struct {
	UUID     string `json:"uuid"`
	Realname string `json:"realname"`
	Tips     string `json:"tips"`
	Status   string `json:"status"`
}

type medialistForJqGrid struct {
	//Rows    []mediaInfoForJqGrid `json:"rows,omitempty"`
	Rows    []interface{} `json:"rows,omitempty"`
	Records int           `json:"records,omitempty"`
	Page    int           `json:"page,omitempty"`
	Total   int           `json:"total,omitempty"`
}

type MediaInfo struct {
	MediaDuration float64 `json:"duration"`
	MediaFrames   int     `json:"frames"`
	MediaFps      float64 `json:"fps"`
	MediaHeight   int     `json:"height"`
	MediaWidth    int     `json:"width"`
	MediaIndex    string  `json:"index"`
	MediaURL      string  `json:"url"`
	MustHLS       bool    `json:"must_hls"`
}

func MediaGetOperation(ctx *gin.Context) {
	mediaIndex := ctx.Param("mediaIndex")
	mediaOp := ctx.Param("mediaOperate")

	info, err := module.MediaGet(mediaIndex)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(42000, err.Error()))
		return
	}

	log("d", "db media type:", info.MediaType)

	switch info.MediaType {

	case global.MediaTypeUltrasonicVideo, global.MediaTypeUltrasonicImage:
		switch mediaOp {
		case "mp4":
			src := path.Join(global.GetPaths().Media, info.Path)
			target := src
			switch path.Ext(src) {
			case ".ogv":
				target = strings.ReplaceAll(src, ".ogv", ".mp4")
			case ".mp4":
				target = src
			default:
				log("e", "unknown source:", src)
				return
			}
			stat, err := os.Stat(target)
			if os.IsNotExist(err) || stat.Size() == 0 {
				log("d", "convert media:", src)
				err = media.ConvertFormat(media.MediaInfo{
					Filepath: src,
				}, media.MediaInfo{
					Filepath:  target,
					MediaType: media.MP4V,
				})
				if err != nil {
					log("e", "convert format:", err.Error())
					return
				}
			}
			ctx.File(target)

		case "ogv":
			if path.Ext(info.Path) == ".ogv" {
				target := path.Join(global.GetPaths().Media, info.Path)
				ctx.File(target)
			} else {
				ctx.JSON(http.StatusOK, FailReturn(402003, "ogv not found"))
			}

		default:
			ctx.JSON(http.StatusOK, FailReturn(402003, fmt.Sprintf("unknow mediaOp %s for ultrasonic video/image", mediaOp)))

		}
		return

	case pacs_dcm4chee.MEDIA_TYPE_IMAGE:
		if len(mediaIndex) != 32 {
			// instance_id
			switch mediaOp {
			case "image":
				bs, _, err := module.PacsGetInstanceMedia(mediaIndex, "image")

				if err != nil {
					ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
				} else {
					ctx.Writer.Write(bs)
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

			case "mp4", "ogv":
				bs, _, err := module.PacsGetInstanceMedia(mediaIndex, mediaOp)

				if err != nil {
					ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
				} else {
					ctx.Writer.Write(bs)
				}

			default:
				ctx.JSON(http.StatusOK, FailReturn(405, fmt.Sprintf("image not fit: %s", mediaOp)))
			}
		} else {
			// media_index

		}

	case pacs_dcm4chee.MEDIA_TYPE_MULTI_FRAME:
		// instance_id

		switch mediaOp {
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

		case "mp4", "ogv":
			bs, _, err := module.PacsGetInstanceMedia(mediaIndex, mediaOp)
			if err != nil {
				ctx.JSON(http.StatusOK, FailReturn(404, err.Error()))
			} else {
				ctx.Writer.Write(bs)
			}

		case "stream", "video.m3u8":
			if path.Ext(info.Path) == ".m3u8" {
				ctx.File(info.Path)
			} else {
				ctx.JSON(http.StatusOK, FailReturn(402003, "unable to get stream"))
			}
			return

		}

	default:
		if path.Ext(info.Path) == ".m3u8" {
			mediaDir := path.Dir(info.Path)
			target := path.Join(mediaDir, mediaOp)
			ctx.File(target)
		}
		return
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
