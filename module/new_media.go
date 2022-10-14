/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: new_media.go
 * Date: 2022/09/23 15:50:23
 */

package module

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"strings"
)

func MediaGetLabelUUIDs(mediaUUID string) (labelIds []string, err error) {
	info, err := MediaGet(mediaUUID)
	if err != nil {
		return nil, err
	}

	for _, labelId := range info.LabelAuthors {
		labelIds = append(labelIds, labelId)
	}
	return labelIds, nil
}

func MediaAddReviewer(mediaUUID string, reviewer int, labelUUID string) error {
	log.Println("i", "media add reviewer", mediaUUID, reviewer, labelUUID)
	info, err := MediaGet(mediaUUID)
	if err != nil {
		return err
	}

	reviewers := info.LabelReviewers
	if uid, ok := reviewers[labelUUID]; !ok || uid < 0 {
		reviewers[labelUUID] = reviewer
		bs, _ := json.Marshal(reviewers)
		return database.GetMedia().LabelUpdateReview(mediaUUID, string(bs))
	} else if uid != reviewer {
		return fmt.Errorf("审阅者UID校验失败：已由他人[%d]审阅", reviewer)
	} else {
		return nil
	}
}

func MediaGetFirstLabeler(mediaUUID string) (labelUUID string, authorUid, reviewerUid int, err error) {
	info, err := MediaGet(mediaUUID)
	if err != nil {
		return "", 0, 0, err
	}

	for uid, label := range info.LabelAuthors {
		authorUid = uid
		labelUUID = label
		break
	}

	reviewerUid = info.LabelReviewers[labelUUID]
	return labelUUID, authorUid, reviewerUid, nil
}

func MediaRemoveLabel(mediaUUID, labelUUID string) error {
	info, err := MediaGet(mediaUUID)
	if err != nil {
		return err
	}

	for uid, uuid := range info.LabelAuthors {
		if uuid == labelUUID {
			delete(info.LabelAuthors, uid)
			bs, _ := json.Marshal(info.LabelAuthors)
			if err = database.GetMedia().LabelUpdateAuthor(mediaUUID, string(bs)); err != nil {
				log.Println("e", err.Error())
			}

			delete(info.LabelReviewers, uuid)
			bs, _ = json.Marshal(info.LabelReviewers)
			if err = database.GetMedia().LabelUpdateAuthor(mediaUUID, string(bs)); err != nil {
				log.Println("e", err.Error())
			}
			break
		}
	}

	database.GetMedia().LabelUpdateProgress(mediaUUID, 0)
	return LabelDelete(labelUUID)
}
func MediaRemoveLabelAll(mediaUUID string) error {
	uuids, err := MediaGetLabelUUIDs(mediaUUID)
	if err != nil {
		return err
	}

	for _, labelUUID := range uuids {
		MediaRemoveLabel(mediaUUID, labelUUID)
	}

	return nil
}

func MediaGetInstanceIdsByStudiesIds(studiesIds []string, ignoreProgress7Check bool) (instanceIds []string, err error) {
	var seriesIds []string

	ps := database.BridgeGetPacsServerHandler()

	for _, studiesId := range studiesIds {
		studiesInfo, err := ps.FindStudiesById(studiesId)
		if err != nil {
			log.Println("e", "find studies id:", err.Error())
			continue
		}

		if ignoreProgress7Check || studiesInfo.LabelProgress == 7 {
			for _, seriesId := range strings.Split(studiesInfo.IncludeSeries, "|") {
				seriesIds = append(seriesIds, seriesId)
			}
		}
	}

	if infos, err := ps.FindSeriesByIdsLocal(seriesIds); err != nil {
		return nil, err
	} else {
		for _, seriesInfo := range infos {
			for _, instanceId := range strings.Split(seriesInfo.IncludeInstances, "|") {
				instanceIds = append(instanceIds, instanceId)
			}
		}
		return instanceIds, nil
	}
}
