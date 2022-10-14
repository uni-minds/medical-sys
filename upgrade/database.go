package upgrade

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"strings"
	"time"
	"uni-minds.com/medical-sys/database"
	"uni-minds.com/medical-sys/global"
	"uni-minds.com/medical-sys/module"
	"uni-minds.com/medical-sys/updater"
)

type UserInfo struct {
	Id            bson.ObjectId `bson:"_id"`
	Username      string        `bson:"username" form:"username"`
	Password      string        `bson:"password" form:"password"`
	Groups        []string      `bson:"groups" form:"groups"`
	Email         string        `bson:"email" form:"email"`
	Realname      string        `bson:"realname" form:"realname"`
	IsExpired     bool
	IsConfirmed   bool
	LoginCount    int `bson:"count"`
	LastPageStart int
}

type MediaInfo struct {
	Mid         bson.ObjectId   `bson:"_id"`
	Filename    string          `bson:"fn"`
	FullPath    string          `bson:"path"`
	Type        string          `bson:"type"`
	OwnerUID    bson.ObjectId   `bson:"ouid"`
	OwnerGroup  string          `bson:"og"`
	LabelUIDs   []bson.ObjectId `bson:"luid"`
	LabelData   []string        `bson:"ld"`
	LabelFrames int             `bson:"lfs"`
	LabelNumber int             `bson:"lnum"`
	Comment     string          `bson:"comment"`
	Height      string          `bson:"h"`
	Width       string          `bson:"w"`
	Frames      string          `bson:"frames"`
	Duration    string          `bson:"duration"`
	UpdateTime  time.Time       `bson:"udt"`
	UploadTime  time.Time       `bson:"ult"`
}

func UpgradeImportUsers() {
	var userList []UserInfo
	updater.ConnectDB()
	_ = updater.GetDB().C("user").Find(nil).All(&userList)

	for _, v := range userList {
		_ = module.UserCreate(v.Username, v.Password, v.Email, v.Realname, v.Id.Hex())
		uid := module.UserGetUid(v.Username)
		for _, g := range v.Groups {
			if g == "admin" {
				g = "administrators"
			}
			gid := module.GroupGetGid(g)
			if gid == 0 {
				module.GroupCreate(g, g, g)
				gid = module.GroupGetGid(g)
				module.GroupAddUser(gid, 1, "master")
			}
			module.GroupAddUser(gid, uid, "member")
		}
	}
	updater.DisconnectDB()
}
func UpgradeImportGroupMedia(group string) {
	var mediaList []MediaInfo
	updater.ConnectDB()
	_ = updater.GetDB().C("media").Find(nil).All(&mediaList)
	gid := module.GroupGetGid(group)
	if gid == 0 {
		log.Panic("No group")
	}

	for _, v := range mediaList {
		if v.OwnerGroup != group {
			continue
		}

		hash := "IMPORT_EMPTY_" + v.Mid.Hex()
		disp := strings.Split(v.Filename, ".")

		f, _ := strconv.Atoi(v.Frames)
		d, _ := strconv.ParseFloat(v.Duration, 32)
		h, _ := strconv.Atoi(v.Height)
		w, _ := strconv.Atoi(v.Width)

		mi := database.MediaInfo{
			Mid:             0,
			DisplayName:     disp[0],
			Path:            v.FullPath,
			Hash:            hash,
			Duration:        d,
			Frames:          f,
			Width:           w,
			Height:          h,
			Status:          0,
			UploadTime:      v.UploadTime.Format(global.TimeFormat),
			UploadUid:       1,
			PatientID:       "",
			MachineID:       "",
			FolderName:      "",
			Fcode:           "",
			IncludeViews:    `["4AP"]`,
			Keywords:        "",
			Memo:            v.Mid.Hex(),
			MediaType:       "",
			MediaData:       "",
			LabelAuthorsUid: "",
			LabelAuthorsLid: "",
			LabelReviewsUid: "",
			LabelReviewsLid: "",
		}

		mid, _ := database.MediaCreate(mi)
		detail := database.MediaInfoUltrasonicVideo{
			PathRaw:  "",
			HashRaw:  "",
			PathJpgs: "",
			Encoder:  "ogv_import",
		}

		err := database.GroupAddMedia(gid, mid)
		if err != nil {
			log.Println(err.Error())
		}
		_ = database.MediaUpdateDetail(mid, detail)

		for i, ld := range v.LabelData {
			u := module.UserGetUidFromMemo(v.LabelUIDs[i].Hex())
			li := database.LabelInfo{
				Uid:        u,
				Mid:        mid,
				Type:       global.LabelTypeAuthor,
				Data:       ld,
				CreateTime: v.UpdateTime.Format(global.TimeFormat),
				ModifyTime: "",
				Memo:       v.LabelUIDs[i].Hex(),
				Frames:     v.LabelFrames,
				Counts:     v.LabelNumber,
				Version:    1,
			}
			database.LabelCreate(li)
			if i > 1 {
				log.Println(mid, module.UserGetRealname(u))
			}
		}
	}
}
func UpgradeParseMediaData() {
	mids := module.GroupGetMedia(2)
	for _, mid := range mids {
		mi, _ := database.MediaGet(mid)
		log.Println(mi.Path)
		p := strings.Replace(mi.Path, `www/files/common/`, `/data/media/medical-sys/us/20200327-12H/`, -1)
		log.Println("=>", p)
		err := database.MediaUpdatePath(mid, p)
		if err != nil {
			fmt.Println(mid, err.Error())
		}
	}
}
