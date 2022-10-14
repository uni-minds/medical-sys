/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: label.go
 */

package module

import (
	"errors"
	"fmt"
	"time"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/global"
)

type LabelSummaryInfo struct {
	AuthorRealname string
	ReviewRealname string
	AuthorProgress string
	ReviewProgress string
	AuthorTips     string
	ReviewTips     string
}
type LabelAuthorInfo struct {
	Realname   string
	Frames     int
	Counts     int
	UpdateTime string
	Progress   string
	Hash       string
	Memo       string
	IsReviewed bool
	IsModified bool
	Json       string
}
type LabelReviewInfo struct {
	Realname   string
	UpdateTime string
	Progress   string
	Tips       string
	Hash       string
	Memo       string
	Author     string
	AuthorTime string
	AuthorJson string
	Json       string
}

//func LabelGetAuthorJson(mid,uid int) (jdata string ,isReviewed,isModifyAfterReview bool,err error) {
//	var li database.LabelInfo
//
//	if uid == -1 {
//		lis, err := database.LabelGetAll(mid, 0, global.LabelTypeAuthor)
//		if err != nil || len(lis) ==0 {
//			return
//		}
//		fmt.Println("MSG Under 多标注未融合")
//		li = lis[0]
//	} else {
//		li, err = database.LabelQuery(mid, uid, global.LabelTypeAuthor)
//		if err != nil {
//			return
//		}
//	}
//
//	if li.Data != "" {
//		var authorData database.LabelInfoAuthorData
//		json.Unmarshal([]byte(li.Data),&authorData)
//		return authorData.Json,authorData.Reviewed,authorData.ModifyAfterReview,nil
//	} else {
//		return
//	}
//}

func LabelGetRealname(i interface{}) (authorName, reviewName string) {
	mi, err := database.MediaGet(i)
	if err != nil {
		log("i", "E LabelGetRealname", err.Error())
		return "", ""
	}
	if mi.LabelAuthorUid > 0 {
		ui, err := database.UserGet(mi.LabelAuthorUid)
		if err != nil {
			log("i", "E LabelGetRealneme", err.Error())
		}
		authorName = ui.Realname
	}
	if mi.LabelReviewUid > 0 {
		ui, err := database.UserGet(mi.LabelReviewUid)
		if err != nil {
			log("i", "E LabelGetRealneme", err.Error())
		}
		reviewName = ui.Realname
	}
	return authorName, reviewName
}

func LabelGetSummary(i interface{}) (summary LabelSummaryInfo, err error) {

	const prog_free = "free"
	const prog_using = "using"
	const prog_submit = "submit"
	const prog_author_reject = "a_reject"
	const prog_review_warning = "r_warning"
	const prog_review_confirm = "r_confirm"

	li, err := database.LabelGet(i)
	if err != nil {
		return
	}
	switch li.Progress {
	case 1:
		// 标注中
		summary.AuthorProgress = prog_using
		summary.ReviewProgress = ""
	case 2:
		// 标注完成
		summary.AuthorProgress = prog_submit
		summary.ReviewProgress = prog_free
	case 3:
		// 审阅中
		summary.AuthorProgress = prog_submit
		summary.ReviewProgress = prog_using
	case 4:
		// 审阅完成，拒绝
		summary.AuthorProgress = prog_author_reject
		summary.ReviewProgress = prog_submit
	case 5:
		// 标注修改中
		summary.AuthorProgress = prog_using
		summary.ReviewProgress = prog_submit
	case 6:
		// 标注完成修改，提交审阅
		summary.AuthorProgress = prog_submit
		summary.ReviewProgress = prog_review_warning
	case 7:
		// 审阅接受，最终状态
		summary.AuthorProgress = prog_submit
		summary.ReviewProgress = prog_review_confirm
	default:
		summary.AuthorProgress = prog_free
		summary.ReviewProgress = ""
	}
	if li.AuthorUid > 0 {
		summary.AuthorRealname = UserGetRealname(li.AuthorUid)
	}
	if li.ReviewUid > 0 {
		summary.ReviewRealname = UserGetRealname(li.ReviewUid)
	}
	return
}
func LabelGetJson(i interface{}) string {
	li, err := database.LabelGet(i)
	if err != nil {
		log("i", "E LabelGetJson", err.Error())
		return ""
	} else {
		return li.Data
	}
}
func LabelUpdateAuthor(jstr string, mediaHash string, uid int) error {
	li, err := database.LabelGet(mediaHash)
	if err != nil {
		// 没有对应的媒体标注，新建
		li = database.LabelInfo{
			Lid:               0,
			Progress:          1,
			AuthorUid:         uid,
			ReviewUid:         0,
			MediaHash:         mediaHash,
			Data:              jstr,
			Version:           1,
			Frames:            0,
			Counts:            0,
			TimeAuthorStart:   time.Now().Format(global.TimeFormat),
			TimeAuthorSubmit:  "",
			TimeReviewStart:   "",
			TimeReviewSubmit:  "",
			TimeReviewConfirm: "",
			Memo:              "",
		}
		return database.LabelCreate(li)

	} else {
		// 存在标注信息，验证是否允许修改
		if li.AuthorUid != uid {
			return errors.New("非原始作者，禁止修改标注")
		}
		switch li.Progress {
		case 3: // fin,ing
			return errors.New("审阅进行中，禁止修改标注")

		case 4: // reject,fin => ing,fin
			li.Progress = 5
		}
		li.Data = jstr
		li.TimeAuthorSubmit = time.Now().Format(global.TimeFormat)
		return database.LabelUpdate(li)
	}
}
func LabelSubmitAuthor(mediaHash string, uid int) error {
	li, err := database.LabelGet(mediaHash)
	if err != nil {
		return err
	}

	if li.AuthorUid != uid {
		return errors.New("非原始作者，禁止提交")
	}

	switch li.Progress {
	case 0, 1: // ing,NaN
		li.Progress = 2
	case 2, 6:
		return errors.New("status conflict: submitted")
	case 3:
		return errors.New("status conflict: under reviewing")
	case 4, 5: // ing,fin
		li.Progress = 6
	case 7:
		return errors.New("status conflict: reviewed")
	default:
		return errors.New("status error")
	}

	li.TimeAuthorSubmit = time.Now().Format(global.TimeFormat)
	return database.LabelUpdate(li)
}
func LabelUpdateReview(jstr string, mediaHash string, uid int) error {
	li, err := database.LabelGet(mediaHash)
	if err != nil {
		return errors.New("无原始数据，无法审阅")

	} else {
		// 存在标注信息，验证是否允许修改
		if li.ReviewUid != uid && li.ReviewUid != 0 {
			ui, _ := database.UserGet(li.ReviewUid)
			return errors.New(fmt.Sprintf("标注已进入专家“%s”的审阅流程，无法修改", ui.Realname))
		}
		li.ReviewUid = uid

		switch li.Progress {
		case 7: // fin,ing
			return errors.New("审阅已经通过，无法修改")

		}
		li.Data = jstr
		li.TimeReviewSubmit = time.Now().Format(global.TimeFormat)
		return database.LabelUpdate(li)
	}
}
func LabelSubmitReview(mediaHash string, uid int, result string) error {
	li, err := database.LabelGet(mediaHash)
	if err != nil {
		log("i", "E1:", err.Error())
		return err
	}

	if li.ReviewUid != uid && li.ReviewUid != 0 {
		return errors.New("非原始审阅者，禁止提交")
	}
	li.ReviewUid = uid

	switch result {
	case "reject":
		switch li.Progress {
		case 2, 3, 6:
			li.Progress = 4
		default:
			return errors.New("审核状态图错误")
		}
		li.TimeReviewSubmit = time.Now().Format(global.TimeFormat)
	case "confirm":
		switch li.Progress {
		case 2, 3, 4, 6:
			li.Progress = 7
		default:
			return errors.New("审核状态图错误")
		}
		li.TimeReviewConfirm = time.Now().Format(global.TimeFormat)
	}
	err = database.LabelUpdate(li)
	if err != nil {
		log("i", "E2:", err.Error())
	}
	return err
}

//func LabelGetSummary(i interface{}, uid int) (realname, activeText, tips string, err error) {
//	li, err := database.LabelGet(i)
//	if err != nil {
//		log("i","LabelGetSummary E", err.Error())
//	}
//	realname = LabelGetRealname(i)
//	if li.Uid == i {
//		activeText = "修改"
//	} else {
//		activeText = "查看"
//	}
//
//	latestTime := li.CreateTime
//	if li.ModifyTime != "" {
//		latestTime = li.ModifyTime
//	}
//
//	switch li.Type {
//	case global.LabelTypeAuthor:
//		tips = fmt.Sprintf("当前进度：%s\n更新时间：%s\n标注帧数：%d\n结构数量：%d\n备注：%s",
//			li.Progress, latestTime, li.Frames, li.Counts, li.Memo)
//
//	case global.LabelTypeReview:
//		tips = fmt.Sprintf("[基于某人的数据]\n当前进度：%s\n更新时间：%s\n标注帧数：%d\n结构数量：%d\n备注：%s",
//			li.Progress, latestTime, li.Frames, li.Counts, li.Memo)
//
//	}
//	return
//}
