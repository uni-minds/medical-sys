/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiLabel.go
 */

package controller

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type LabelData struct {
	MediaHash string `json:"media"`
	Data      string `json:"data"`
	Direction string `json:"direction"`
	Admin     string `json:"admin"`
}

type LabelInfoForButton struct {
	TextBackG string `json:"textbackg"`
	TextHover string `json:"texthover"`
	Tips      string `json:"tips"`
}

func LabelGet(ctx *gin.Context) {
	uid := -1
	if uidi, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = uidi.(int)
	}

	mediaIndex := ctx.Param("mediaIndex")
	op := ctx.Param("op")
	//mediaType := mediaIndexAnalysis(mediaIndex)

	switch op {
	case "getrealname":
		authorName, reviewName := module.LabelGetRealname(mediaIndex)
		ctx.JSON(http.StatusOK, SuccessReturn([]string{authorName, reviewName}))

	case "summary":
		summary, authorUid, reviewUid, err := module.LabelGetSummary(mediaIndex)
		if err != nil {
			ctx.JSON(http.StatusOK, FailReturn(30001, err.Error()))
		} else if authorUid != uid && reviewUid != uid {
			ctx.JSON(http.StatusOK, FailReturn(http.StatusForbidden, summary))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(summary))
		}

	case "memo":
		memo := module.LabelGetMemo(mediaIndex)
		ctx.JSON(http.StatusOK, SuccessReturn(memo))
		return

	case "author", "review":
		ld := module.LabelGetJson(mediaIndex)
		ctx.JSON(http.StatusOK, SuccessReturn(ld))
		return
	}
}

func LabelPost(ctx *gin.Context) {
	uid := -1
	if uidi, exists := ctx.Get("uid"); !exists {
		return
	} else {
		uid = uidi.(int)
	}

	mediaIndex := ctx.Param("mediaIndex")
	//mtype := mediaIndexAnalysis(mediaIndex)

	var ldata LabelData
	if err := ctx.BindJSON(&ldata); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return
	}

	switch ctx.Param("op") {
	case "memo":
		if err := module.LabelUpdateMemo(mediaIndex, uid, ldata.Data); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		} else {
			ctx.JSON(http.StatusOK, SuccessReturn(fmt.Sprintf("同步成功 @ %s", time.Now().Format(global.TimeFormat))))
		}
		return

	case "author":
		switch ctx.Query("do") {
		case "full":
			if err := module.LabelUpdateAuthor(ldata.Data, mediaIndex, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				log("i", fmt.Sprintf("label update by author: %s", module.UserGetRealname(uid)))
				ctx.JSON(http.StatusOK, SuccessReturn(fmt.Sprintf("同步成功 @ %s", time.Now().Format(global.TimeFormat))))
			}

		case "submit":
			// 提交
			if err := module.LabelSubmitAuthor(mediaIndex, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn("exit"))
			}
			return
		}

	case "review":
		switch ctx.Query("do") {
		case "full":
			if err := module.LabelUpdateReview(ldata.Data, mediaIndex, uid); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				return

			} else {
				fmt.Println("review submit:", module.UserGetRealname(uid))
				ctx.JSON(http.StatusOK, SuccessReturn(fmt.Sprintf("同步成功 @ %s", time.Now().Format(global.TimeFormat))))
				return
			}

		case "reject":
			if err := module.LabelSubmitReview(mediaIndex, uid, "reject"); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				return
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				return
			}

		case "confirm":
			if err := module.LabelSubmitReview(mediaIndex, uid, "confirm"); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				return
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				return
			}

		case "revoke":
			if ldata.Admin == global.DefAdminPassword {
				fmt.Println("review revoke")
				if err := module.LabelRevokeReview(mediaIndex, 0, true); err != nil {
					ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
				} else {
					fmt.Println("review revoke by", module.UserGetRealname(uid))
					ctx.JSON(http.StatusOK, SuccessReturn("exit"))
				}
			}
		}
	}
}

func LabelDelete(ctx *gin.Context) {
	mediaIndex := ctx.Param("mediaIndex")

	var ldata LabelData
	if err := ctx.BindJSON(&ldata); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return
	}

	switch ctx.Param("op") {
	case "full":
		if ldata.Admin == global.DefAdminPassword {
			if err := module.LabelDelete(mediaIndex); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
			} else {
				ctx.JSON(http.StatusOK, SuccessReturn("exit"))
			}
		} else {
			fmt.Println("admin@az")
			ctx.JSON(http.StatusOK, FailReturn(http.StatusForbidden, "forbidden"))
		}
	}
}

func mediaIndexAnalysis(mediaIndex string) (mediaType string) {
	if strings.Contains(mediaIndex, DICOM_TYPE_US_ID) {
		return "dicom"
	} else if len(mediaIndex) == 32 {
		return "hash"
	} else {
		return ""
	}
}
