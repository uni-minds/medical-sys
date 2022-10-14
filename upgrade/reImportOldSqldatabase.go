/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  reImportOldSqldatabase
 * @Version: 1.0.0
 * @Date: 2020/4/18 14:02
 */

package upgrade

import (
	"fmt"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	"uni-minds.com/liuxy/medical-sys/database"
)

var once sync.Once
var engin *gorose.Engin
var usersDB map[string]int

type labelDetail struct {
	Disp  string `json:"Disp"`
	User  string `json:"User"`
	Ctime string `json:"ct"`
	Data  string `json:"dt"`
}

var labels1DB, labels2DB map[string]labelDetail

type labelInfo struct {
	Lid   int
	Uid   int
	Mid   int
	Type  string
	Data  string
	CTime string `gorose:"createtime"`
	MTime string `gorose:"modifytime"`
}

func (*labelInfo) TableName() string {
	return "label"
}

type mediaInfo struct {
	Mid      int
	DispName string `gorose:"displayname"`
	Hash     string
	Memo     string
}

func (*mediaInfo) TableName() string {
	return "media"
}

type userInfo struct {
	Uid      int
	Username string
	Groups   string
}

func (*userInfo) TableName() string {
	return "users"
}

func Run3() {
	lis, err := database.LabelGetAll()
	if err != nil {
		panic(err)
	}
	for _, li := range lis {
		mi, err := database.MediaGet(li.MediaHash)
		if err != nil {
			fmt.Println(err.Error(), li.MediaHash)
			continue
		}
		fmt.Println(li.Progress, li.AuthorUid, li.ReviewUid)
		database.MediaUpdateLabelProgress(mi.Mid, li.AuthorUid, li.ReviewUid, li.Progress)
	}
}
