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
	"gitee.com/uni-minds/utils/tools"
	"strconv"
	"strings"
)

const (
	GroupTypeLabelMedia        = "label_media"
	GroupTypeLabelDicom        = "label_dicom"
	GroupTypeLabelStream       = "label_stream"
	GroupTypeScreenDicom       = "screen_dicom"
	GroupTypeAdmin             = "admin"
	GroupContainTypeUser       = "user"
	GroupContainTypeMediaId    = "media_uuid"
	GroupContainTypeInstanceId = "instance_id"
	GroupContainTypeHls        = "hls"
	GroupContainTypeStudiesId  = "studies_id"
)

func initGroupDB() {
	dbSql := fmt.Sprintf(`create table IF NOT EXISTS "%s" (
    id             INTEGER           not null
        primary key autoincrement,
    name           TEXT default ''   not null,
    type           TEXT default ''   not null,
    users          TEXT default '{}' not null,
    display_name   TEXT default ''   not null,
    contain_data   TEXT default ''   not null,
    contain_type   TEXT default ''   not null,
    contain_counts INT  default 0    not null,
    memo           TEXT default ''   not null
);`, global.DefaultDatabaseGroupTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Error(err.Error())
	}
}

func GroupGetAll() (groupInfos []DbStructGroup, err error) {
	err = DB().Table(&groupInfos).Select()
	return
}

func GroupCreate(info DbStructGroup) (gid int, err error) {
	gt, err := GroupGet(info.Id)
	if err != nil {
		info.Id = 0
		info.Name = strings.ToLower(info.Name)
		_, err = DB().Table(global.DefaultDatabaseGroupTable).Data(info).Insert()
		info, _ = GroupGet(info.Name)
		return info.Id, err
	}

	return gt.Id, errors.New(global.EGroupAlreadyExisted)
}
func GroupGet(i interface{}) (info DbStructGroup, err error) {
	switch i.(type) {
	case int:
		err = DB().Table(&info).Where("id", "=", i).Select()

	case string:
		err = DB().Table(&info).Where("name", "=", strings.ToLower(i.(string))).Select()
	}

	if err != nil {
		return info, err
	} else if info.Id == 0 {
		return info, fmt.Errorf(global.EGroupNotExisted)
	} else {
		return info, nil
	}
}

func GroupDelete(i interface{}) error {
	switch i.(type) {
	case int:
		if i.(int) > 1 {
			_, err := DB().Table(global.DefaultDatabaseGroupTable).Where("id", i).Delete()
			return err
		} else {
			return fmt.Errorf("组序号异常: %d", i)
		}

	case string:
		if i.(string) != "" {
			_, err := DB().Table(global.DefaultDatabaseGroupTable).Where("name", i).Delete()
			return err
		} else {
			return fmt.Errorf("组序号异常: %d", i)
		}

	default:
		return errors.New(global.EGroupForbidden)
	}
}

func groupUpdateCore(id int, data interface{}) error {
	_, err := DB().Table(global.DefaultDatabaseGroupTable).Data(data).Where("id", id).Update()
	return err
}

func GroupUpdateName(gid int, gn string) error {
	data := map[string]interface{}{"name": gn}
	return groupUpdateCore(gid, data)
}
func GroupUpdateDisplayName(gid int, name string) error {
	data := map[string]interface{}{"display_name": name}
	return groupUpdateCore(gid, data)
}
func GroupUpdateGroupType(gid int, tp string) error {
	data := map[string]interface{}{"type": tp}
	return groupUpdateCore(gid, data)
}
func GroupUpdateContainType(gid int, tp string) error {
	data := map[string]interface{}{"contain_type": tp}
	return groupUpdateCore(gid, data)
}

func GroupGetContainId(id int) (mids []int, err error) {
	mIndex, mType, err := GroupGetContains(id)
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

func GroupGetContains(id int) (containIndex []string, containType string, err error) {
	info, err := GroupGet(id)
	if err != nil {
		return nil, "", err
	}

	var data []string
	if err := json.Unmarshal([]byte(info.ContainData), &data); err != nil {
		data = make([]string, 0)

	}

	switch info.Type {
	case GroupTypeLabelDicom, GroupTypeScreenDicom, GroupTypeLabelMedia, GroupTypeLabelStream, "default":
		return data, info.ContainType, err

	case GroupTypeAdmin:
		return nil, "admin", nil

	default:
		return nil, "", fmt.Errorf("unknown group type: %s", info.Type)

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

	data := map[string]interface{}{"contain_counts": len(index), "contain_data": string(bs)}
	return groupUpdateCore(gid, data)
}

func GroupRemoveMedia(gid, mid int) error {
	mids, err := GroupGetContainId(gid)
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
		p, err := GroupGetUserPermissions(gl.Id, uid)
		if err != nil || !p.ListMedia {
			continue
		} else {
			gids = append(gids, gl.Id)
		}
	}
	return gids, nil
}
