/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: users.go
 */

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/tools"
	"sort"
	"strconv"
	"strings"
	"time"
)

type UserListContent struct {
	Uid         int    `json:"uid"`
	Username    string `json:"username"`
	Realname    string `json:"realname"`
	Email       string `json:"email"`
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
func UserGetGroups(uid int, groupType string) (gids []int) {
	gids = make([]int, 0)
	groupType = strings.ToLower(groupType)
	if UserIsAdmin(uid) {
		gl, _ := database.GroupGetAll()
		for _, v := range gl {
			if groupType == "*" || strings.Contains(v.GroupType, groupType) {
				gids = append(gids, v.Gid)
			}
		}
	} else {
		userGroupIds, err := database.GroupListByUser(uid)
		if err != nil {
			log("e", err.Error())
		}

		for _, gid := range userGroupIds {
			if gi, err := database.GroupGet(gid); err != nil {
				log("e", err.Error())
			} else if groupType == "*" || strings.Contains(gi.GroupType, groupType) {
				gids = append(gids, gid)
			}
		}
	}
	return gids
}

func UserGetGroupContains(uid, gid int) (mediaIndex []string, mediaType string, err error) {
	level := GroupGetUserLevel(gid, uid)
	if level < GroupLevelGuest && !UserIsAdmin(uid) {
		return nil, "", errors.New("user forbidden")
	}

	return database.GroupGetContains(gid)
}

func UserGetGroupContainsSelector(uid, gid int, sortField, sortOrder string) []string {
	mediaIndex, _, err := UserGetGroupContains(uid, gid)
	if err != nil {
		return nil
	}

	// prepare database
	data := make(map[string]interface{}, 0)
	switch sortField {
	case "view":
		for _, mid := range mediaIndex {
			id, err := strconv.Atoi(mid)
			if err != nil {
				// mediahash
				mi, _ := database.MediaGet(mid)
				data[mid] = mi.IncludeViews
			} else {
				// mid
				mi, _ := database.MediaGet(id)
				data[mid] = mi.IncludeViews
			}
		}

	case "memo":
		for _, mid := range mediaIndex {
			id, err := strconv.Atoi(mid)
			if err != nil {
				// mediahash
				mi, _ := database.MediaGet(mid)
				data[mid] = mi.Memo
			} else {
				mi, _ := database.MediaGet(id)
				data[mid] = mi.Memo
			}
		}

	case "duration":
		for _, mid := range mediaIndex {
			id, err := strconv.Atoi(mid)
			if err != nil {
				// mediahash
				mi, _ := database.MediaGet(mid)
				data[mid] = mi.Duration
			} else {
				mi, _ := database.MediaGet(id)
				data[mid] = mi.Duration
			}
		}

	case "name":
		for _, mid := range mediaIndex {
			id, err := strconv.Atoi(mid)
			if err != nil {
				// mediahash
				mi, _ := database.MediaGet(mid)
				data[mid] = mi.DisplayName
			} else {
				mi, _ := database.MediaGet(id)
				data[mid] = mi.DisplayName
			}
		}

	case "frames":
		for _, mid := range mediaIndex {
			id, err := strconv.Atoi(mid)
			if err != nil {
				// mediahash
				mi, _ := database.MediaGet(mid)
				data[mid] = mi.Frames
			} else {
				mi, _ := database.MediaGet(id)
				data[mid] = mi.Frames
			}
		}

	case "authors":
		for _, mid := range mediaIndex {
			id, err := strconv.Atoi(mid)
			if err != nil {
				// mediahash
				mi, _ := database.MediaGet(mid)
				data[mid] = mi.LabelAuthorUid
			} else {
				mi, _ := database.MediaGet(id)
				data[mid] = mi.LabelAuthorUid
			}
		}

	case "reviews":
		for _, mid := range mediaIndex {
			id, err := strconv.Atoi(mid)
			if err != nil {
				// mediahash
				mi, _ := database.MediaGet(mid)
				data[mid] = mi.LabelReviewUid
			} else {
				mi, _ := database.MediaGet(id)
				data[mid] = mi.LabelReviewUid
			}
		}

	default:
		// mid
		switch sortOrder {
		case "desc":
			sort.Sort(sort.Reverse(sort.StringSlice(mediaIndex)))

		default:
			sort.Strings(mediaIndex)
		}
		return mediaIndex
	}

	//d := tools.MediaSorter(data)
	//switch sortOrder {
	//case "desc":
	//	// 降序
	//	sort.Sort(sort.Reverse(d))
	//
	//default:
	//	// 升序
	//	sort.Sort(d)
	//}
	//
	//mediaIndex,mediaType,err = make([]int, 0)
	//for _, v := range d {
	//	mediaIndex,mediaType,err = append(mediaIndex,mediaType,err, v.Mid)
	//}
	return mediaIndex
}

func UserGetUidFromMemo(memo string) int {
	u, err := database.UserGetManual("memo", memo)
	if err != nil {
		log("i", "Find by memo E:", err.Error())
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
		log("i", "E UserGetMid", err.Error())
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
		log("i", "Admins override.")
		return tmp, nil
	}

	permissionViewMedia := false

	gids := UserGetGroups(uid, "")
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
	log("i", "Memo update:", mid, uid, mi.Memo, "->", memo)
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
		log("i", v)
		//g := UserGroupInfo{
		//	Gid:       v.Gid,
		//	GroupName: v.GroupName,
		//}
	}
	return
}

// passwords
func UserCheckPassword(username, password string) (uid int) {
	u, err := database.UserGet(username)
	if err != nil {
		log("e", "userLogin:", err.Error())
		return -1 //用户不存在
	}
	if u.Activate == 0 {
		log("e", fmt.Sprintf("login user: %s is deactivate", username))
		return -2 //账号未确认
	}

	result := tools.GetStringMD5(password + u.PasswordSalt)
	if u.LoginFailCount > 1 {
		time.Sleep(3 * time.Second)
	}

	if u.Password != result {
		t := u.LoginFailCount + 1
		log("e", fmt.Sprintf("login failed: %s, %d times", username, t))
		database.UserUpdateTryFailureCount(u.Uid, t)
		if t >= 5 {
			log("e", "user account now deactivate")
			database.UserUpdateAccountActiveType(u.Uid, -1)
		}
		return -3 //密码错
	}

	if u.ExpireTime != "" {
		t, err := time.ParseInLocation(global.TimeFormat, u.ExpireTime, time.Local)
		if err == nil && t.Sub(time.Now()) <= 0 {
			return -4 // 账号过期
		}
	}

	log("i", fmt.Sprintf("login success: %s", username))
	database.UserUpdateTryFailureCount(u.Uid, 0)
	database.UserUpdateLoginCount(u.Uid, u.LoginCount+1)

	return u.Uid
}
func UserSetPassword(i interface{}, password string) error {
	u, err := database.UserGet(i)
	if err != nil {
		return err
	}

	salt := tools.RandStringFromAlphabet(20, "")
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
func UserGetAll() map[int]database.UserInfo {

	uis, err := database.UserGetAll()
	if err != nil {
		log("e", "userGetAll:", err.Error())
		return nil
	}

	data := make(map[int]database.UserInfo, 0)
	for _, v := range uis {
		data[v.Uid] = v
	}
	return data
}
