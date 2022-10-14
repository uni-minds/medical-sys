/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: group.go
 */

package module

import (
	"errors"
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

func GroupCreate(groupname, displayname, groupType, containType, memo string) error {
	gi := database.GroupInfo{
		GroupName:   groupname,
		DisplayName: displayname,
		GroupType:   groupType,
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
	return gi.Gid
}

func GroupGetType(gid int) string {
	gi, err := database.GroupGet(gid)
	if err != nil {
		return ""
	} else {
		return gi.GroupType
	}
}
func GroupGetGroupname(gid int) string {
	gi, err := database.GroupGet(gid)
	if err != nil {
		return ""
	}
	return gi.GroupName
}
func GroupGetDisplayname(gid int) string {
	gi, err := database.GroupGet(gid)
	if err != nil {
		return ""
	}
	return gi.DisplayName
}
func GroupGetMedia(gid int) (media []int) {
	media, err := database.GroupGetMediaId(gid)
	if err != nil {
		return nil
	}
	return
}
func GroupGetDicom(gid int) (studiesIds []string) {
	studiesIds, _, err := database.GroupGetContains(gid)
	if err != nil {
		return nil
	}
	return
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
	database.GroupDelete(gi.Gid)
	return nil
}

func GroupGetAll() map[int][]string {
	gis, err := database.GroupGetAll()
	if err != nil {
		log("e", "GroupGetAll:", err.Error())
		return nil
	}

	data := make(map[int][]string, 0)
	for _, v := range gis {
		data[v.Gid] = []string{v.GroupName, v.DisplayName, v.Users}
	}
	return data
}
