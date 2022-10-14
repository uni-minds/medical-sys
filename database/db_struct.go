/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: db_struct.go
 * Date: 2022/09/27 09:00:27
 */

package database

import "gitee.com/uni-minds/medical-sys/global"

type DbStructGroup struct {
	Id            int    `gorose:"id"`
	Name          string `gorose:"name"`
	Type          string `gorose:"group_type"`
	Users         string `gorose:"users"`
	ContainType   string `gorose:"contain_type"`
	ContainCounts int    `gorose:"contain_counts"`
	ContainData   string `gorose:"contain_data"`
	Memo          string `gorose:"memo"`
}

func (*DbStructGroup) TableName() string {
	return global.DefaultDatabaseGroupTable
}

type DbStructUser struct {
	Uid      int    `gorose:"uid"`
	Username string `gorose:"username"`
	Realname string `gorose:"realname"`
	//Groups         string `gorose:"groups"`
	Password       string `gorose:"password"`
	PasswordSalt   string `gorose:"passwordsalt"`
	Email          string `gorose:"email"`
	Activate       int    `gorose:"activate"`
	ExpireTime     string `gorose:"expiretime"`
	Memo           string `gorose:"memo"`
	RegisterTime   string `gorose:"registertime"`
	LoginCount     int    `gorose:"logincount"`
	LoginTime      string `gorose:"logintime"`
	LoginFailCount int    `gorose:"loginfailcount"`
	LastGroupId    int    `gorose:"lastGroupId"`
	LastPageIndex  int    `gorose:"lastPageIndex"`
	LastToken      string `gorose:"lastToken"`
}

func (*DbStructUser) TableName() string {
	return global.DefaultDatabaseUserTable
}

type DbStructLabel struct {
	Id                int    `gorose:"id"`
	LabelUUID         string `gorose:"label_uuid"`
	MediaUUID         string `gorose:"media_uuid"`
	Progress          int    `gorose:"progress"`
	AuthorUid         int    `gorose:"author"`
	ReviewUid         int    `gorose:"reviewer"`
	Data              string `gorose:"data"`
	Version           int    `gorose:"version"`
	Frames            int    `gorose:"frames"`
	Counts            int    `gorose:"counts"`
	TimeAuthorStart   int64  `gorose:"t_author_start"`
	TimeAuthorSubmit  int64  `gorose:"t_author_submit"`
	TimeReviewStart   int64  `gorose:"t_reviewer_start"`
	TimeReviewSubmit  int64  `gorose:"t_reviewer_submit"`
	TimeReviewConfirm int64  `gorose:"t_reviewer_confirm"`
	Memo              string `gorose:"memo"`
}

func (*DbStructLabel) TableName() string {
	return global.DefaultDatabaseLabelTable
}

type DbStructMedia struct {
	Id            int     `gorose:"id"`
	MediaUUID     string  `gorose:"media_uuid"`
	DisplayName   string  `gorose:"display"`
	Path          string  `gorose:"path"`
	Width         int     `gorose:"width"`
	Height        int     `gorose:"height"`
	Duration      float64 `gorose:"duration"`
	Frames        int     `gorose:"frames"`
	Fps           float64 `gorose:"fps"`
	MediaType     string  `gorose:"media_type"`
	MediaHash     string  `gorose:"media_hash"`
	UploadUid     int     `gorose:"upload_uid"`
	UploadTime    int64   `gorose:"upload_time"`
	Memo          string  `gorose:"memo"`
	PatientId     string  `gorose:"patient_id"`
	MachineId     string  `gorose:"machine_id"`
	Metadata      string  `gorose:"metadata"`
	CrfDefine     string  `gorose:"crf_define"`
	Keywords      string  `gorose:"keywords"`
	MediaData     string  `gorose:"media_data"`
	LabelAuthor   string  `gorose:"label_author"`
	LabelReviewer string  `gorose:"label_reviewer"`
	LabelProgress int     `gorose:"label_progress"`
	CoworkType    string  `gorose:"cowork"`
}

func (*DbStructMedia) TableName() string {
	return global.DefaultDatabaseMediaTable
}
