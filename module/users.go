package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
	"uni-minds.com/medical-sys/database"
	"uni-minds.com/medical-sys/global"
	"uni-minds.com/medical-sys/tools"
)

type UserListContent struct {
	Uid         int    `json:"uid"`
	Username    string `json:"username"`
	Realname    string `json:"realname"`
	Email       string `json:"email"`
	Groups      string `json:"groups"`
	LoginCount  int    `json:"logincount"`
	LoginTime   string `json:"logintime"`
	LoginEnable bool   `json:"loginenable"`
	Remark      string `json:"remark"`
}

type UserGroupInfo struct {
	Gid       int
	GroupName string
	Role      string
	Memo      string
}

func UserCreate(username, password, email, realname, memo string) error {
	u := database.UserInfo{
		Username: username,
		Email:    email,
		Realname: realname,
		Memo:     memo,
	}
	uid, err := database.UserCreate(u)
	if err != nil {
		return err
	}
	return UserSetPassword(uid, password)
}
func UserSetActive(uid int) error {
	return database.UserUpdateAccountActiveType(uid, 1)
}
func UserGetUid(username string) int {
	u, err := database.UserGet(username)
	if err != nil {
		return 0
	} else {
		return u.Uid
	}
}
func UserGetRealname(i interface{}) string {
	u, err := database.UserGet(i)
	if err != nil {
		return ""
	}
	return u.Realname
}
func UserGetGroups(uid int) (gids []int) {
	gids = make([]int, 0)
	if UserIsAdmin(uid) {
		gl, _ := database.GroupGetAll()
		for _, v := range gl {
			gids = append(gids, v.Gid)
		}
	} else {
		gids, _ = database.UserGetGroups(uid)
	}
	return gids
}
func UserGetGroupMedia(uid, gid int) (mids []int) {
	level := GroupGetUserLevel(gid, uid)
	if level < GroupLevelGuest && !UserIsAdmin(uid) {
		return nil
	}

	mids, err := database.GroupGetMedia(gid)
	if err != nil {
		log.Println(err.Error())
	}
	return mids
}
func UserGetGroupMediaSelector(uid, gid int, sortField, sortOrder string) []int {
	mids := UserGetGroupMedia(uid, gid)

	// prepare database
	data := make(map[int]interface{}, 0)
	switch sortField {
	case "view":
		for _, mid := range mids {
			mi, _ := database.MediaGet(mid)
			data[mid] = mi.IncludeViews
		}

	case "memo":
		for _, mid := range mids {
			mi, _ := database.MediaGet(mid)
			data[mid] = mi.Memo
		}

	case "duration":
		for _, mid := range mids {
			mi, _ := database.MediaGet(mid)
			data[mid] = mi.Duration
		}

	case "name":
		for _, mid := range mids {
			mi, _ := database.MediaGet(mid)
			data[mid] = mi.DisplayName
		}

	case "frames":
		for _, mid := range mids {
			mi, _ := database.MediaGet(mid)
			data[mid] = mi.Frames
		}

	case "authors":
		for _, mid := range mids {
			mi, _ := database.MediaGet(mid)
			data[mid] = mi.LabelAuthorsUid
		}

	case "reviews":
		for _, mid := range mids {
			mi, _ := database.MediaGet(mid)
			data[mid] = mi.LabelReviewsUid
		}

	default:
		// mid
		switch sortOrder {
		case "desc":
			sort.Sort(sort.Reverse(sort.IntSlice(mids)))

		default:
			sort.Ints(mids)
		}
		return mids
	}

	d := tools.MediaSorter(data)
	switch sortOrder {
	case "desc":
		// 降序
		sort.Sort(sort.Reverse(d))

	default:
		// 升序
		sort.Sort(d)
	}

	mids = make([]int, 0)
	for _, v := range d {
		mids = append(mids, v.Mid)
	}
	return mids
}
func UserGetUidFromMemo(memo string) int {
	u, err := database.UserGetManual("memo", memo)
	if err != nil {
		log.Println("Find by memo E:", err.Error())
		return 0
	}
	return u.Uid
}
func UserGetMediaHash(uid, mid int) string {
	mi, err := userGetMediaInfo(uid, mid)
	if err != nil {
		return ""
	} else {
		return mi.Hash
	}
}
func UserGetMediaMemo(uid, mid int) string {
	mi, err := userGetMediaInfo(uid, mid)
	if err != nil {
		return ""
	} else {
		return mi.Memo
	}
}
func UserGetMid(uid int, hash string) int {
	mi, err := userGetMediaInfo(uid, hash)
	if err != nil {
		log.Println("E UserGetMid", err.Error())
		return 0
	} else {
		return mi.Mid
	}
}
func userGetMediaInfo(uid int, i interface{}) (mi database.MediaInfo, err error) {
	tmp, err := database.MediaGet(i)
	if err != nil {
		return
	}
	if UserIsAdmin(uid) {
		log.Println("Admins override.")
		return tmp, nil
	}

	permissionViewMedia := false

	gids := UserGetGroups(uid)
	for _, gid := range gids {
		mids := GroupGetMedia(gid)
		for _, mid := range mids {
			if mid == tmp.Mid {
				if GroupGetUserLevel(gid, uid) > GroupLevelNotAllowed {
					permissionViewMedia = true
					break
				}
			}
		}
		if permissionViewMedia {
			break
		}
	}
	if permissionViewMedia {
		return tmp, nil
	} else {
		return mi, errors.New(global.EMediaForbidden)
	}
}

func UserIsAdmin(uid int) bool {
	gi, err := database.GroupGet(global.DefaultAdminGroup)
	if err != nil {
		return false
	}

	users, err := database.GroupGetUsers(gi.Gid)
	_, ok := users[uid]
	return ok
}

func UserSetMediaMemo(uid, mid int, memo string) error {
	mi, err := userGetMediaInfo(uid, mid)
	if err != nil {
		return err
	}
	log.Println("Memo update:", mid, uid, mi.Memo, "->", memo)
	return database.MediaUpdateMemo(mid, memo)
}

func UserList() (jsonstr string) {
	userlist := make([]UserListContent, 0)
	db, err := database.UserGetAll()
	if err != nil {
		return
	}

	for _, v := range db {
		e := UserListContent{
			Uid:         v.Uid,
			Username:    v.Username,
			Realname:    v.Realname,
			Email:       v.Email,
			Groups:      "",
			LoginEnable: v.Activate == 1,
			LoginCount:  v.LoginCount,
			LoginTime:   strings.Replace(v.LoginTime, "T", "\n", 1),
			Remark:      v.Memo,
		}
		userlist = append(userlist, e)
	}

	d, err := json.Marshal(userlist)
	if err != nil {
		return
	}
	return string(d)
}

func UserGroupsList(uid int) (groups []UserGroupInfo) {
	if uid <= 0 {
		return
	}

	gl, err := database.GroupGetAll()
	if err != nil {
		return
	}
	groups = make([]UserGroupInfo, 0)

	for _, v := range gl {
		log.Println(v)
		//g := UserGroupInfo{
		//	Gid:       v.Gid,
		//	GroupName: v.GroupName,
		//}
	}
	return
}

// passwords
func UserCheckPassword(username, password string) (uid int) {
	fmt.Println("checking:", username, password)
	u, err := database.UserGet(username)
	if err != nil {
		log.Println(err.Error())
		return -1 //用户不存在
	}
	fmt.Println("calc target", u.Password, u.PasswordSalt)
	if u.Activate == 0 {
		return -2 //账号未确认
	}

	result := tools.GetStringMD5(password + u.PasswordSalt)
	fmt.Println("calc result", result)
	if u.Password != result {
		t := u.LoginFailCount + 1
		database.UserUpdateTryFailureCount(u.Uid, t)
		if t >= 5 {
			database.UserUpdateAccountActiveType(u.Uid, -1)
		}
		log.Println("Login failure count", t)
		return -3 //密码错
	}

	if u.ExpireTime != "" {
		t, err := time.ParseInLocation(global.TimeFormat, u.ExpireTime, time.Local)
		if err == nil && t.Sub(time.Now()) <= 0 {
			return -4 // 账号过期
		}
	}

	database.UserUpdateTryFailureCount(u.Uid, 0)
	database.UserUpdateLoginCount(u.Uid, u.LoginCount+1)

	return u.Uid
}
func UserSetPassword(i interface{}, password string) error {
	u, err := database.UserGet(i)
	if err != nil {
		return err
	}

	salt := tools.GenSaltString(20)
	password = tools.GetStringMD5(password + salt)
	return database.UserUpdatePassword(u.Uid, password, salt)
}
func UserChangePassword(username, oldPassword, newPassword string) (ok bool) {
	uid := UserCheckPassword(username, oldPassword)
	if uid > 0 {
		return nil == UserSetPassword(username, newPassword)
	}
	return false
}

func UserDelete(username string) error {
	u, err := database.UserGet(username)
	if err != nil {
		return err
	}
	return database.UserDelete(u.Uid)
}
