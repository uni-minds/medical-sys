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
	pacs_global "gitee.com/uni-minds/bridge-pacs/global"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/utils/tools"
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
	u := database.DbStructUser{
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
			if groupType == "*" || strings.Contains(v.Type, groupType) {
				gids = append(gids, v.Id)
			}
		}
	} else {
		userGroupIds, err := database.GroupListByUser(uid)
		if err != nil {
			log.Println("e", err.Error())
		}

		for _, gid := range userGroupIds {
			if gi, err := database.GroupGet(gid); err != nil {
				log.Println("e", err.Error())
			} else if groupType == "*" || strings.Contains(gi.Type, groupType) {
				gids = append(gids, gid)
			}
		}
	}
	return gids
}

func UserGetGroupContains(uid, gid int) (idx []string, containType string, err error) {
	level := GroupGetUserLevel(gid, uid)
	if UserIsAdmin(uid) || global.FlagGetDebug() || level >= GroupLevelGuest {
		return database.GroupGetContains(gid)
	} else {
		return nil, "", errors.New("user forbidden")
	}
}

func UserGetGroupContainsWithSelector(uid, gid int, selectView, sortField, sortOrder string, start, count int) (result []MediaInfo, total int, err error) {
	mediaUUIDs, mediaType, err := UserGetGroupContains(uid, gid)
	if err != nil {
		return nil, 0, err
	}

	switch mediaType {
	case database.GroupContainTypeHls, database.GroupContainTypeMediaId:
		infos, err := database.GetMedia().Selector(mediaUUIDs, sortField, sortOrder == "asc")
		if err != nil {
			return nil, 0, err
		} else {
			total = len(infos)

			if len(infos) > start {
				infos = infos[start:]
			} else {
				return nil, total, nil
			}

			if len(infos) > count {
				infos = infos[:count]
			}

			result = make([]MediaInfo, 0)
			for _, info := range infos {
				if info.Height == 0 || info.Width == 0 || info.Fps <= 0 {
					if err = MediaRescan(info.MediaUUID); err != nil {
						log.Println("e", "media rescan", info.MediaUUID, err.Error())
						continue
					}

					info, _ = database.GetMedia().Get(info.MediaUUID)
				}
				result = append(result, ConvertMediaInfoFromDbMedia(info))
			}
			return result, total, nil
		}

	case database.GroupContainTypeStudiesId:
		instanceIDs, err := MediaGetInstanceIdsByStudiesIds(mediaUUIDs, true)
		if err != nil {
			return nil, 0, err
		}

		ps := database.BridgeGetPacsServerHandler()
		infos, err := ps.FindInstanceByIdsLocal(instanceIDs, 10000)
		if err != nil {
			return nil, 0, err
		}

		// ignore instance without label view
		tmp := make([]pacs_global.InstanceInfo, 0)
		for _, info := range infos {
			if "" != info.LabelView && ("" == selectView || info.LabelView == selectView) {
				tmp = append(tmp, info)
			}
		}

		total = len(tmp)

		if len(tmp) > start {
			tmp = tmp[start:]
		} else {
			return nil, total, nil
		}

		if len(tmp) > count {
			tmp = tmp[:count]
		}

		result := ConvertMediaInfoFromInstanceInfos(tmp)
		return result, total, nil

	default:
		return nil, 0, fmt.Errorf("unknow media type: %s", mediaType)
	}
}

func UserGetUidFromMemo(memo string) int {
	u, err := database.UserGetManual("memo", memo)
	if err != nil {
		log.Println("i", "Find by memo E:", err.Error())
		return 0
	}
	return u.Uid
}

func UserGetMediaHash(uid, mid int) string {
	mi, err := userGetMediaInfo(uid, mid)
	if err != nil {
		return ""
	} else {
		return mi.MediaUUID
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
		log.Error(fmt.Sprintf("UserGetMid: %s", err.Error()))
		return 0
	} else {
		return mi.Id
	}
}
func userGetMediaInfo(uid int, i interface{}) (mi MediaInfo, err error) {
	info, err := MediaGet(i)
	if err != nil {
		return
	}
	if UserIsAdmin(uid) {
		log.Println("Admins override.")
		return info, nil
	}

	permissionViewMedia := false

	gids := UserGetGroups(uid, "")
	for _, gid := range gids {
		containIndex, _ := GroupGetContainMedia(gid)
		for _, mediaUUID := range containIndex {
			if mediaUUID == info.MediaUUID {
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
		return info, nil
	} else {
		return mi, errors.New(global.EMediaForbidden)
	}
}

func UserIsAdmin(uid int) bool {
	gi, err := database.GroupGet(global.DefaultAdminGroup)
	if err != nil {
		return false
	}

	users, err := database.GroupGetUsers(gi.Id)
	_, ok := users[uid]
	return ok
}

func UserSetMediaMemo(uid, mid int, memo string) error {
	mi, err := userGetMediaInfo(uid, mid)
	if err != nil {
		return err
	}
	log.Println("i", "Memo update:", mid, uid, mi.Memo, "->", memo)
	return database.GetMedia().UpdateMemo(mid, memo)
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
		log.Println("i", v)
		//g := UserGroupInfo{
		//	Id:       v.Id,
		//	Name: v.Name,
		//}
	}
	return
}

// passwords
func UserCheckPassword(username, password string) (uid int) {
	u, err := database.UserGet(username)
	if err != nil {
		log.Println("e", "userLogin:", err.Error())
		return -1 //用户不存在
	}
	if u.Activate == 0 {
		log.Println("e", fmt.Sprintf("login user: %s is deactivate", username))
		return -2 //账号未确认
	}

	result := tools.GetStringChecksum(password+u.PasswordSalt, tools.ModeChecksumMD5)
	if u.LoginFailCount > 1 {
		time.Sleep(3 * time.Second)
	}

	if u.Password != result {
		t := u.LoginFailCount + 1
		log.Println("e", fmt.Sprintf("login failed: %s, %d times", username, t))
		database.UserUpdateTryFailureCount(u.Uid, t)
		if t >= 5 {
			log.Println("e", "user account now deactivate")
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

	log.Println("i", fmt.Sprintf("login success: %s", username))
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
	password = tools.GetStringChecksum(password+salt, tools.ModeChecksumMD5)
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
func UserGetAll() map[int]database.DbStructUser {

	uis, err := database.UserGetAll()
	if err != nil {
		log.Println("e", "userGetAll:", err.Error())
		return nil
	}

	data := make(map[int]database.DbStructUser, 0)
	for _, v := range uis {
		data[v.Uid] = v
	}
	return data
}
