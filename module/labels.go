package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"uni-minds.com/medical-sys/database"
	"uni-minds.com/medical-sys/global"
)

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
func LabelGetAuthorJson(i interface{}) string {
	var li database.LabelInfo
	var err error
	li, err = database.LabelGet(i)
	if err != nil {
		log.Println("LabelGetAuthorJson", err.Error())
		return ""
	}
	switch li.Type {
	case global.LabelTypeAuthor:
		var authorData database.LabelInfoAuthorData
		json.Unmarshal([]byte(li.Data), &authorData)
		return authorData.Json

	}
	return ""
}
func LabelGetReviewJson(i interface{}) (reviewJson string, authorJson string, err error) {
	var li database.LabelInfo
	li, err = database.LabelGet(i)
	if err != nil {
		log.Println("LabelGetReviewJson", err.Error())
		return "", "", err
	}
	switch li.Type {
	case global.LabelTypeReview:
		var reviewData database.LabelInfoReviewerData
		err = json.Unmarshal([]byte(li.Data), &reviewData)
		return reviewData.Json, reviewData.BasedJson, err
	}
	return
}
func LabelSetReviewerJson(mid, uid int, json string) error {
	if mid <= 0 || uid <= 0 {
		return errors.New("Invalid mid/uid")
	}

	f, c := LabelParseReviewJson(json, 1)
	li, err := database.LabelQuery(mid, uid, global.LabelTypeReview)
	if err != nil {
		li := database.LabelInfo{
			Lid:      0,
			Uid:      uid,
			Mid:      mid,
			Version:  1,
			Progress: global.LabelProgressReviewing,
			Type:     global.LabelTypeReview,
			Data:     json,
			Frames:   f,
			Counts:   c,
		}
		_, _, err := database.LabelCreate(li)
		return err
	} else {
		revdata := database.LabelInfoReviewerData{
			Json: json,
		}
		err = database.LabelUpdateLabelData(li.Lid, uid, f, c, revdata, global.LabelProgressReviewing)
		return err
	}
}
func LabelSetAuthorJson(mid, uid int, json string) error {
	if mid <= 0 || uid <= 0 {
		return errors.New("Invalid mid/uid")
	}
	f, c := LabelParseReviewJson(json, 1)
	li, err := database.LabelQuery(mid, uid, global.LabelTypeAuthor)
	if err != nil {
		li = database.LabelInfo{
			Lid:      0,
			Uid:      uid,
			Mid:      mid,
			Version:  1,
			Progress: global.LabelProgressAuthoring,
			Frames:   f,
			Counts:   c,
			Type:     global.LabelTypeAuthor,
		}
		li.Lid, li.Hash, err = database.LabelCreate(li)
		if err != nil {
			log.Println("E LabelSetAuthorJSON", err.Error())
		}
	}
	authordata := database.LabelInfoAuthorData{
		Json:              json,
		Reviewed:          false,
		ModifyAfterReview: false,
	}
	return database.LabelUpdateLabelData(li.Lid, uid, f, c, authordata, global.LabelTypeAuthor)
}

//func LabelCreate(mid,uid int,)
func LabelGetHashs(lids []int) (hashs []string) {
	hashs = make([]string, 0)
	for _, lid := range lids {
		li, err := database.LabelGet(lid)
		if err != nil {
			log.Println("labelget error", err.Error())
			continue
		}
		hashs = append(hashs, li.Hash)
	}
	return
}
func LabelGetRealname(i interface{}) string {
	li, err := database.LabelGet(i)
	if err != nil {
		log.Println("LabelGetRealname E", err.Error())
		return ""
	} else {
		ui, err := database.UserGet(li.Uid)
		if err != nil {
			log.Println("LabelGetRealneme E", err.Error())
			return ""
		} else {
			return ui.Realname
		}
	}
}
func LabelGetSummary(i interface{}, uid int) (realname, activeText, tips string, err error) {
	li, err := database.LabelGet(i)
	if err != nil {
		log.Println("LabelGetSummary E", err.Error())
	}
	realname = LabelGetRealname(i)
	if li.Uid == i {
		activeText = "修改"
	} else {
		activeText = "查看"
	}

	latestTime := li.CreateTime
	if li.ModifyTime != "" {
		latestTime = li.ModifyTime
	}

	switch li.Type {
	case global.LabelTypeAuthor:
		tips = fmt.Sprintf("当前进度：%s\n更新时间：%s\n标注帧数：%d\n结构数量：%d\n备注：%s",
			li.Progress, latestTime, li.Frames, li.Counts, li.Memo)

	case global.LabelTypeReview:
		tips = fmt.Sprintf("[基于某人的数据]\n当前进度：%s\n更新时间：%s\n标注帧数：%d\n结构数量：%d\n备注：%s",
			li.Progress, latestTime, li.Frames, li.Counts, li.Memo)

	}
	return
}
