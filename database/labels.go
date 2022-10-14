/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  labels
 * @Version: 1.0.0
 * @Date: 2020/4/14 23:50
 */

package database

import (
	"errors"
	"fmt"
	"log"
	"uni-minds.com/medical-sys/global"
)

func (*LabelsInfo) TableName() string {
	return global.DefaultDatabaseLabelsTableNew
}

type LabelsInfo struct {
	Lid               int    `gorose:"lid"`
	Progress          int    `gorose:"progress"`
	AuthorUid         int    `gorose:"authorUid"`
	ReviewUid         int    `gorose:"reviewUid"`
	MediaHash         string `gorose:"mediaHash"`
	Data              string `gorose:"data"`
	Version           int    `gorose:"version"`
	Frames            int    `gorose:"frames"`
	Counts            int    `gorose:"counts"`
	TimeAuthorStart   string `gorose:"timeAuthorStart"`
	TimeAuthorSubmit  string `gorose:"timeAuthorSubmit"`
	TimeReviewStart   string `gorose:"timeReviewStart"`
	TimeReviewSubmit  string `gorose:"timeReviewSubmit"`
	TimeReviewConfirm string `gorose:"timeReviewConfirm"`
	Memo              string `gorose:"memo"`
}

func LabelsGet(i interface{}) (li LabelsInfo, err error) {
	switch i.(type) {
	case int:
		err = DB().Table(&li).Where("lid", "=", i).Select()
		if err != nil || li.Lid == 0 {
			err = errors.New(global.ELabelDBLabedNotExist)
		}
	case string:
		err = DB().Table(&li).Where("mediahash", "=", i).Select()
		if err != nil || li.Lid == 0 {
			err = errors.New(global.ELabelDBLabedNotExist)
		}
	}
	return
}
func LabelsCreate(li LabelsInfo) error {
	//li.Lid=0
	_, err := DB().Table(global.DefaultDatabaseLabelsTableNew).Data(li).Insert()
	if err != nil {
		log.Println("E Label create:", err.Error())
	}
	return err
}
func LabelsUpdate(li LabelsInfo) (err error) {
	_, err = DB().Table(global.DefaultDatabaseLabelsTableNew).Data(li).Where("lid", "=", li.Lid).Update()
	if err != nil {
		fmt.Println("DB E", err.Error())
	}
	return
}
