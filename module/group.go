/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: group.go
 */

package module

import (
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
)

const (
	GroupLevelCustom     = -1
	GroupLevelNotAllowed = 0
	GroupLevelGuest      = 1
	GroupLevelMember     = 2
	GroupLevelLeader     = 3
	GroupLevelMaster     = 4
)

func GroupCreate(groupName, groupType, containType, memo string) error {
	gi := database.DbStructGroup{
		Name:        groupName,
		Type:        groupType,
		ContainType: containType,
		Memo:        memo,
	}
	_, err := database.GroupCreate(gi)
	return err
}

func GroupGetGid(groupname string) int {
	gi, err := database.GroupGet(groupname)
	if err != nil {
		return 0
	}
	return gi.Id
}

func GroupGetType(gid int) string {
	gi, err := database.GroupGet(gid)
	if err != nil {
		return ""
	} else {
		return gi.Type
	}
}
func GroupGetGroupname(gid int) string {
	gi, err := database.GroupGet(gid)
	if err != nil {
		return ""
	}
	return gi.Name
}
func GroupGetDisplayname(gid int) string {
	gi, err := database.GroupGet(gid)
	if err != nil {
		return ""
	}
	return gi.Name
}
func GroupGetContains(id int) (containIndex []string, containType string, err error) {
	return database.GroupGetContains(id)
}

func GroupAddMedia(gid int, ids []string) error {
	log.Trace(fmt.Sprintf("group %d add media %v", gid, ids))
	return database.GroupAddContains(gid, ids)
}

func GroupGetContainMedia(id int) (containIndex []string, containType string) {
	var err error
	containIndex, containType, err = database.GroupGetContains(id)
	if err != nil {
		log.Println("e", "group get contain media:", err.Error())
		return nil, ""
	} else {
		return containIndex, containType
	}
}
func GroupGetContainViews(gid, uid int, ignoreProgressCheck bool) (view []string, err error) {
	mediaIndex, mediaType, err := UserGetGroupContains(uid, gid)
	if err != nil {
		return nil, err
	}

	viewMap := make(map[string]struct{}, 0)
	switch mediaType {
	case database.GroupContainTypeStudiesId:
		ps := database.BridgeGetPacsServerHandler()

		if instanceIds, err := MediaGetInstanceIdsByStudiesIds(mediaIndex, ignoreProgressCheck); err != nil {
			return nil, err

		} else if infos, err := ps.FindInstanceByIdsLocal(instanceIds, 0); err != nil {
			return nil, err

		} else {
			for _, info := range infos {
				if info.LabelView == "" {
					// 无切面标注媒体直接忽略
					continue
				}

				if _, ok := viewMap[info.LabelView]; !ok {
					viewMap[info.LabelView] = struct{}{}
				}
			}
		}
	}

	var views []string
	for key, _ := range viewMap {
		views = append(views, key)
	}
	return views, nil
}

func GroupGetUserLevel(gid, uid int) (level int) {
	p, err := database.GroupGetUserPermissions(gid, uid)
	if err != nil {
		return
	}
	if p.ListMedia {
		level = GroupLevelGuest
	}
	if p.ListLabels {
		level = GroupLevelMember
	}
	if p.ManageReviews {
		level = GroupLevelLeader
	}
	if p.ManageUsers {
		level = GroupLevelMaster
	}
	return
}

func GroupUserAdd(gid int, uid int, role string) error {
	p := database.UserPermissions{}

	p.ListMedia = true

	if role == "member" || role == "leader" || role == "master" {
		p.ListLabels = true
		p.ListUsers = true
	}

	if role == "leader" || role == "master" {
		p.ManageReviews = true
		p.ListReviews = true
	}

	if role == "master" {
		p.ManageLabels = true
		p.ManageReviews = true
		p.ManageUsers = true
	}

	return database.GroupAddUser(gid, uid, p)
}
func GroupUserAddFrendly(groupname, username, role string) error {
	gid := GroupGetGid(groupname)
	uid := UserGetUid(username)
	return GroupUserAdd(gid, uid, role)
}
func GroupUserSetPermissioin(gid, uid int, role string) error {
	permissions := database.UserPermissions{}
	switch role {
	case "guest":
		permissions.ListMedia = true
	case "member":
		permissions.ListLabels = true
	case "leader":
		permissions.ManageReviews = true
	case "master":
		permissions.ManageUsers = true
	default:
		return errors.New("unknow permission:" + role)
	}
	return database.GroupSetUserPermissions(gid, uid, permissions)
}
func GroupDel(i interface{}) error {
	gi, err := database.GroupGet(i)
	if err != nil {
		return err
	}
	database.GroupDelete(gi.Id)
	return nil
}

func GroupGetAll() map[int][]string {
	gis, err := database.GroupGetAll()
	if err != nil {
		log.Println("e", "GroupGetAll:", err.Error())
		return nil
	}

	data := make(map[int][]string, 0)
	for _, v := range gis {
		data[v.Id] = []string{v.Name, v.Name, v.Users}
	}
	return data
}
