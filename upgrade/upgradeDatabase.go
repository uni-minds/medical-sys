package upgrade

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"strings"
	"time"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/module"
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
	ConnectDB()
	_ = GetDB().C("User").Find(nil).All(&userList)

	for _, v := range userList {
		if module.UserGetUid(v.Username) > 0 {
			continue
		}
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
				module.GroupUserAdd(gid, 1, "master")
			}
			module.GroupUserAdd(gid, uid, "member")
		}
	}
	DisconnectDB()
}
func UpgradeImportGroupMedia(group string) {
	var userList []UserInfo
	ConnectDB()
	_ = GetDB().C("User").Find(nil).All(&userList)

	userdb := make(map[string]string, 0)

	for _, v := range userList {
		userdb[v.Id.Hex()] = v.Username
	}

	var mediaList []MediaInfo
	ConnectDB()
	_ = GetDB().C("media").Find(nil).All(&mediaList)
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

		if disp[0] == "腹动F35764-483650 (1)" {
			log.Println("Catch")
		}

		f, _ := strconv.Atoi(v.Frames)
		d, _ := strconv.ParseFloat(v.Duration, 32)
		h, _ := strconv.Atoi(v.Height)
		w, _ := strconv.Atoi(v.Width)

		s1 := strings.Replace(v.FullPath, `www/files/G2020C_AC/`, `/Data/media/medical-sys/us/20200415-08H/`, -1)
		s2 := strings.Replace(s1, `www/files/G2020B_AC/`, `/Data/media/medical-sys/us/20200415-08H/`, -1)
		s3 := strings.Replace(s2, `www/files/G2020B_AC/`, `/Data/media/medical-sys/us/20200415-08H/`, -1)

		mi := database.MediaInfo{
			Mid:          0,
			DisplayName:  disp[0],
			Path:         s3,
			Hash:         hash,
			Duration:     d,
			Frames:       f,
			Width:        w,
			Height:       h,
			Status:       0,
			UploadTime:   v.UploadTime.Format(global.TimeFormat),
			UploadUid:    1,
			PatientID:    "",
			MachineID:    "",
			FolderName:   "",
			Fcode:        "",
			IncludeViews: `["A"]`,
			Keywords:     `["正常"]`,
			Memo:         v.Mid.Hex(),
			MediaType:    "",
			MediaData:    "",
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
			u, err := database.UserGet(userdb[v.LabelUIDs[i].Hex()])
			if err != nil {
				log.Println("User Get err", err, v.LabelUIDs[i])
				continue
			}
			li := database.LabelInfo{
				Progress:          2,
				AuthorUid:         u.Uid,
				ReviewUid:         0,
				MediaHash:         hash,
				Data:              ld,
				Version:           1,
				Frames:            v.LabelFrames,
				Counts:            v.LabelNumber,
				TimeAuthorStart:   v.UpdateTime.Format(global.TimeFormat),
				TimeAuthorSubmit:  "",
				TimeReviewStart:   "",
				TimeReviewSubmit:  "",
				TimeReviewConfirm: "",
				Memo:              v.LabelUIDs[i].Hex(),
			}
			database.LabelCreate(li)
			if i > 1 {
				log.Println(mid, module.UserGetRealname(u))
			}
		}
	}
}
func UpgradeParseMediaData() {
	mids := module.GroupGetMedia(11)
	for _, mid := range mids {
		mi, _ := database.MediaGet(mid)
		log.Println(mi.Path)
		p := strings.Replace(mi.Path, `www/files/G2020A_AC/`, `/Data/media/medical-sys/us/20200415-08H/`, -1)
		log.Println("=>", p)
		err := database.MediaUpdatePath(mid, p)
		if err != nil {
			fmt.Println(mid, err.Error())
		}
	}
}
