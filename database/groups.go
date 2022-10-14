/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: groups_common.go
 */

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/tools"
	"strconv"
	"strings"
)

const (
	GroupTypeLabelMedia  = "label_media"
	GroupTypeLabelDicom  = "label_dicom"
	GroupTypeScreenDicom = "screen_dicom"

	GroupContainTypeInstanceId = "instance_id"
	GroupContainTypeStudiesId  = "studies_id"
	GroupContainTypeMediaId    = "mid"
	GroupContainTypeUser       = "user"
)

type GroupInfo struct {
	Gid         int    `gorose:"gid"`
	GroupName   string `gorose:"group_name"`
	GroupType   string `gorose:"group_type"`
	DisplayName string `gorose:"display_name"`
	Users       string `gorose:"users"`
	Memo        string `gorose:"memo"`
	ContainData string `gorose:"contain_data"`
	ContainType string `gorose:"contain_type"`
	MediaCounts int    `gorose:"media_counts"`
}

func (*GroupInfo) TableName() string {
	return global.DefaultDatabaseGroupTable
}

func initGroupDB() {
	dbSql := fmt.Sprintf(`create table IF NOT EXISTS "%s" (
    gid          INTEGER not null primary key autoincrement,
    group_name   TEXT default "" not null,
    display_name TEXT default "" not null,
    contain_data TEXT default "[]" not null,
    contain_type text,
    users        TEXT default "{}" not null,
    memo         TEXT default "" not null,
    group_type   TEXT default "" not null,
    media_counts int  default 0 not null);`, global.DefaultDatabaseGroupTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log("E", err.Error())
	}
}

func GroupGetAll() (gl []GroupInfo, err error) {
	err = DB().Table(&gl).OrderBy("gid").Select()
	return
}

func GroupCreate(gi GroupInfo) (gid int, err error) {
	gt, err := GroupGet(gi.GroupName)
	if err != nil {
		gi.Gid = 0
		gi.GroupName = strings.ToLower(gi.GroupName)
		_, err = DB().Table(global.DefaultDatabaseGroupTable).Data(gi).Insert()
		gi, _ = GroupGet(gi.GroupName)
		return gi.Gid, err
	}

	return gt.Gid, errors.New(global.EGroupAlreadyExisted)
}
func GroupGet(i interface{}) (gi GroupInfo, err error) {
	switch i.(type) {
	case int:
		err = DB().Table(&gi).Where("gid", "=", i).Select()
		if gi.Gid == 0 {
			err = errors.New(global.EGroupNotExisted)
		}

	case string:
		err = DB().Table(&gi).Where("group_name", "=", strings.ToLower(i.(string))).Select()
		if gi.Gid == 0 {
			err = errors.New(global.EGroupNotExisted)
		}
	}
	return
}
func GroupDelete(gid int) error {
	if gid > 1 {
		_, err := DB().Table(global.DefaultDatabaseGroupTable).Where("gid", "=", gid).Delete()
		return err
	} else {
		return errors.New(global.EGroupForbidden)
	}
}

func groupUpdateCore(gid int, data interface{}) error {
	_, err := DB().Table(global.DefaultDatabaseGroupTable).Data(data).Where("gid", "=", gid).Update()
	return err
}

func GroupUpdateName(gid int, gn string) error {
	data := map[string]interface{}{"group_name": gn}
	return groupUpdateCore(gid, data)
}
func GroupUpdateDisplayName(gid int, name string) error {
	data := map[string]interface{}{"display_name": name}
	return groupUpdateCore(gid, data)
}
func GroupUpdateGroupType(gid int, tp string) error {
	data := map[string]interface{}{"group_type": tp}
	return groupUpdateCore(gid, data)
}
func GroupUpdateContainType(gid int, tp string) error {
	data := map[string]interface{}{"contain_type": tp}
	return groupUpdateCore(gid, data)
}

func GroupGetMediaId(gid int) (mids []int, err error) {
	mIndex, mType, err := GroupGetContains(gid)
	switch mType {
	case GroupContainTypeMediaId:
		for _, index := range mIndex {
			id, _ := strconv.Atoi(index)
			mids = append(mids, id)
		}
		return mids, err
	default:
		return nil, errors.New("group type mismatch")
	}
}

func GroupGetContains(gid int) (mediaIndex []string, mediaType string, err error) {
	gi, err := GroupGet(gid)
	if err != nil {
		return nil, "", err
	}

	var data []string
	err = json.Unmarshal([]byte(gi.ContainData), &data)

	switch gi.GroupType {
	case GroupTypeLabelDicom, GroupTypeScreenDicom:
		return data, gi.ContainType, err

	case GroupTypeLabelMedia:
		return data, gi.ContainType, err

	case "admin":
		return nil, "admin", nil

	default:
		return nil, "", errors.New(fmt.Sprint("unknown group type:", gi.GroupType))

	}
}

func GroupGetUsers(gid int) (users map[int]int, err error) {
	gi, err := GroupGet(gid)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(gi.Users), &users)
	return
}

func GroupAddContain(gid int, id interface{}) (err error) {
	var idx string
	switch id.(type) {
	case string:
		idx = id.(string)
	case int:
		idx = strconv.Itoa(id.(int))
	}

	mediaIndex, _, err := GroupGetContains(gid)
	if err != nil {
		return err
	}

	for _, v := range mediaIndex {
		if id == v {
			return errors.New(global.EGroupMediaAlreadyInThisGroup)
		}
	}

	mediaIndex = append(mediaIndex, idx)
	return groupUpdateContains(gid, mediaIndex)
}

func GroupAddContains(gid int, ids []string) error {
	containIds, _, err := GroupGetContains(gid)
	if err != nil {
		return err
	}

	containIds = tools.StringsDedup(append(containIds, ids...))

	return groupUpdateContains(gid, containIds)
}

func GroupHasContain(gid int, id string) bool {
	containIds, _, err := GroupGetContains(gid)
	if err != nil {
		return false
	}

	for _, hasId := range containIds {
		if hasId == id {
			return true
		}
	}
	return false
}

func GroupAddUser(gid, uid int, permissions UserPermissions) error {
	users, err := GroupGetUsers(gid)
	if err != nil {
		return err
	}

	users[uid] = setPermissions(permissions)
	fmt.Printf("GID= %d add UID= %d...", gid, uid)
	if err = groupUpdateUsers(gid, users); err != nil {
		fmt.Println("E", err.Error())
	} else {
		fmt.Println("OK")
	}
	return err
}
func GroupRemoveUser(gid, uid int) error {
	users, err := GroupGetUsers(gid)
	if err != nil {
		return err
	}

	_, ok := users[uid]
	if !ok {
		fmt.Println("W group don't have user UID=", uid)
		return nil
	}

	delete(users, uid)
	return groupUpdateUsers(gid, users)
}
func groupUpdateUsers(gid int, users map[int]int) error {
	j, err := json.Marshal(users)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"users": string(j)}
	return groupUpdateCore(gid, data)
}

func groupUpdateContains(gid int, mediaIndex interface{}) error {
	index := make([]string, 0)

	switch mediaIndex.(type) {
	case []int:
		for _, v := range mediaIndex.([]int) {
			index = append(index, strconv.Itoa(v))
		}

	case []string:
		index = mediaIndex.([]string)

	default:
		return errors.New("unknown media index type")
	}

	bs, _ := json.Marshal(index)

	data := map[string]interface{}{"media_counts": len(index), "contain_data": string(bs)}
	return groupUpdateCore(gid, data)
}

func GroupRemoveMedia(gid, mid int) error {
	mids, err := GroupGetMediaId(gid)
	if err != nil {
		return err
	}

	for i, v := range mids {
		if v == mid {
			mids = append(mids[:i], mids[i+1:]...)
			return groupUpdateContains(gid, mids)
		}
	}
	return errors.New(global.EGroupMediaNotExistInGroup)
}

func GroupGetUserPermissions(gid, uid int) (permissions UserPermissions, err error) {
	users, err := GroupGetUsers(gid)

	p, ok := users[uid]
	if ok {
		permissions = getPermissions(p)
	} else {
		err = errors.New(global.EGroupUserNotExist)
	}
	return
}
func GroupSetUserPermissions(gid, uid int, permissions interface{}) error {
	users, err := GroupGetUsers(gid)
	if err != nil {
		return err
	}

	_, ok := users[uid]
	if !ok {
		return errors.New(global.EGroupUserNotExist)
	}

	switch permissions.(type) {
	case UserPermissions:
		users[uid] = setPermissions(permissions.(UserPermissions))
		return groupUpdateUsers(gid, users)

	case string:
		users[uid] = setPermissionsRole(permissions.(string))
		return groupUpdateUsers(gid, users)

	default:
		return errors.New("E unknown permission type")
	}
}

func GroupListByUser(uid int) (gids []int, err error) {
	gls, err := GroupGetAll()
	gids = make([]int, 0)
	if err != nil {
		return nil, err
	}

	for _, gl := range gls {
		p, err := GroupGetUserPermissions(gl.Gid, uid)
		if err != nil || !p.ListMedia {
			continue
		} else {
			gids = append(gids, gl.Gid)
		}
	}
	return gids, nil
}
