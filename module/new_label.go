/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: new_label.go
 * Date: 2022/09/23 15:49:23
 */

package module

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"time"
)

func LabelGet(i interface{}) (database.DbStructLabel, error) {
	db := database.GetLabel()
	info, err := db.Get(i)
	if err == nil {
		return info, nil
	}

	// instanceid则提取首个标注
	switch i.(type) {
	case string:
		if global.IsDicomUUID(i.(string)) {
			infos, err := LabelGetByMediaUUID(i.(string))
			if err != nil || len(infos) == 0 {
				return database.DbStructLabel{}, err
			} else {
				return infos[0], nil
			}
		}
	}
	return database.DbStructLabel{}, fmt.Errorf("LabelGet无法解析interface:%v", i)
}

func LabelGetByMediaUUID(mediaUUID string) (infos []database.DbStructLabel, err error) {
	dbInfos, err := database.GetLabel().GetByMediaUUID(mediaUUID)
	if err != nil {
		return nil, err
	} else {
		for _, info := range dbInfos {
			infos = append(infos, info)
		}
		return infos, nil
	}
}

func ProcessAuthorData(mediaUUID string, uid int, data string) error {
	info, err := MediaGet(mediaUUID)
	if err != nil {
		return err
	}

	authors := info.LabelAuthors

	if labelUUID, ok := authors[uid]; ok {
		// 更新Label
		return LabelUpdateAuthor(labelUUID, uid, data)

	} else if len(authors) == 0 || info.CoworkType != global.CoworkTypeSingle {
		// 无标注人或允许并行
		if labelUUID, err = LabelCreate(mediaUUID, uid, data); err != nil {
			// 创建Label
			return err
		}

		log.Println("i", "LabelAddAuthor, media=", mediaUUID, "author=", uid, "label=", labelUUID)

		if global.IsDicomUUID(mediaUUID) {
			log.Println("d", fmt.Sprintf("create_label.instance_id"))
			return nil

		} else {
			log.Println("d", fmt.Sprintf("create_label.update_media info"))
			authors = make(map[int]string)
			authors[uid] = labelUUID
			bs, _ := json.Marshal(authors)
			return database.GetMedia().LabelUpdateAuthor(mediaUUID, string(bs))

		}
	} else {
		return fmt.Errorf("标注作者校验失败：%d | %v", uid, authors)

	}
}

func ProcessAuthorSubmit(mediaUUID string, uid int, labelUUID string) error {
	mi, err := MediaGet(mediaUUID)
	if err != nil {
		return err
	}
	if labelUUID == "" {
		var ok bool
		labelUUID, ok = mi.LabelAuthors[uid]
		if !ok {
			return fmt.Errorf("数据异常：未找到用户对应的标注")
		}
	}

	var labelInfo database.DbStructLabel
	if global.IsDicomUUID(mediaUUID) {
		// 处理Instance媒体

		if labels, err := LabelGetByMediaUUID(mediaUUID); err != nil {
			return err
		} else if len(labels) == 0 {
			return fmt.Errorf("异常：无媒体数据无关联标注，请先添加标注再提交")
		} else {
			labelInfo = labels[0]
		}

	} else {
		// 处理一般媒体信息
		labelInfo, err = LabelGet(labelUUID)
		if err != nil {
			return err
		}
	}

	var nextProgress int
	switch labelInfo.Progress {
	case 1: // ing,NaN
		// 首次提交
		nextProgress = 2

		// 非Dicom数据需要在Media库中添加Reviewer
		if !global.IsDicomUUID(mediaUUID) {
			if err = MediaAddReviewer(mediaUUID, -1, labelUUID); err != nil {
				return err
			}
		}
	case 4, 5: // ing,fin
		// 驳回后提交
		nextProgress = 6
	default:
		return fmt.Errorf("标注状态异常：%s", ProgressQuery(labelInfo.Progress))
	}

	// 更新Label库
	labelInfo.Progress = nextProgress
	labelInfo.TimeAuthorSubmit = time.Now().Unix()

	if nextProgress > mi.LabelProgress {
		MediaSetLabelProgress(mediaUUID, nextProgress)
	}

	return database.GetLabel().UpdateAll(labelInfo)
}

func ProcessReviewConfirm(mediaUUID, labelUUID string, uid int) error {
	mi, err := MediaGet(mediaUUID)
	if err != nil {
		return err
	}

	reviewerUid, ok := mi.LabelReviewers[labelUUID]
	if ok && reviewerUid > 0 && reviewerUid != uid {
		return fmt.Errorf("非原始审阅者(%s)，禁止提交，", UserGetRealname(reviewerUid))
	} else if !ok || reviewerUid < 1 {
		MediaAddReviewer(mediaUUID, uid, labelUUID)
	}

	labelInfo, err := LabelGet(labelUUID)
	if err != nil {
		return err
	}

	const nextProgress = 7

	switch labelInfo.Progress {
	case 2, 3, 5, 4, 6, 7:
		labelInfo.Progress = nextProgress
	default:
		return fmt.Errorf("标注状态异常：%s", ProgressQuery(labelInfo.Progress))
	}

	labelInfo.ReviewUid = uid
	labelInfo.TimeReviewConfirm = time.Now().Unix()

	if mi.LabelProgress != nextProgress {
		MediaSetLabelProgress(mediaUUID, nextProgress)
	}

	return database.GetLabel().UpdateAll(labelInfo)
}

func ProcessReviewReject(mediaUUID, labelUUID string, uid int) error {
	mi, err := MediaGet(mediaUUID)
	if err != nil {
		return err
	}

	reviewerUid, ok := mi.LabelReviewers[labelUUID]
	if ok && reviewerUid > 0 && reviewerUid != uid {
		return fmt.Errorf("非原始审阅者(%s)，禁止提交，", UserGetRealname(reviewerUid))
	} else if !ok || reviewerUid < 1 {
		MediaAddReviewer(mediaUUID, uid, labelUUID)
	}

	labelInfo, err := LabelGet(labelUUID)
	if err != nil {
		return err
	}

	const nextProgress = 4
	switch labelInfo.Progress {
	case 2, 3, 6:
		labelInfo.Progress = nextProgress
	default:
		return fmt.Errorf("标注状态异常：%s", ProgressQuery(labelInfo.Progress))
	}

	labelInfo.ReviewUid = uid
	labelInfo.TimeReviewConfirm = time.Now().Unix()

	if mi.LabelProgress != nextProgress {
		MediaSetLabelProgress(mediaUUID, nextProgress)
	}

	return database.GetLabel().UpdateAll(labelInfo)
}

func ProcessReviewRevoke(mediaUUID, labelUUID string, uid int, superpass string) error {
	log.Println("w", fmt.Sprintf("user %s review revoke", UserGetRealname(uid)))
	if superpass == global.DefAdminPassword {
		mi, err := MediaGet(mediaUUID)
		if err != nil {
			return err
		}

		_, ok := mi.LabelReviewers[labelUUID]
		if !ok {
			return fmt.Errorf("不存在审阅记录")
		}

		labelInfo, err := LabelGet(labelUUID)
		if err != nil {
			return err
		}

		const nextProgress = 2 // 至提交审阅
		labelInfo.Progress = nextProgress

		labelInfo.ReviewUid = -1
		labelInfo.TimeReviewSubmit = time.Now().Unix()

		MediaSetLabelProgress(mediaUUID, nextProgress)
		return database.GetLabel().UpdateAll(labelInfo)

	} else {
		return fmt.Errorf("提权密码错误，暂不支持专家手动撤回审阅")
	}
}
