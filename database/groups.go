/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: groups.go
 */

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"uni-minds.com/liuxy/medical-sys/global"
)

const (
	GroupUserPermissionsListMedia     = 1 << 0
	GroupUserPermissionsManageMedia   = 1 << 1
	GroupUserPermissionsListUsers     = 1 << 2
	GroupUserPermissionsManageUsers   = 1 << 3
	GroupUserPermissionsListLabels    = 1 << 4
	GroupUserPermissionsManageLabels  = 1 << 5
	GroupUserPermissionsListReviews   = 1 << 6
	GroupUserPermissionsManageReviews = 1 << 7
)

type GroupInfo struct {
	Gid         int    `gorose:"gid"`
	GroupName   string `gorose:"groupname"`
	DisplayName string `gorose:"displayname"`
	Media       string `gorose:"media"`
	Users       string `gorose:"users"`
	Memo        string `gorose:"memo"`
}

func initGroupDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"gid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"groupname" TEXT NOT NULL default "",
	"displayname" TEXT NOT NULL default "",	
	"media" TEXT NOT NULL default "[]",
	"users" TEXT NOT NULL default "{}",
	"memo" TEXT NOT NULL default "")`, global.DefaultDatabaseGroupTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Panic(err.Error())
	}
}
func setPermissions(permissions UserPermissions) (p int) {
	if permissions.ListMedia {
		p |= GroupUserPermissionsListMedia
	}

	if permissions.ManageMedia {
		p |= GroupUserPermissionsManageMedia
	}

	if permissions.ListUsers {
		p |= GroupUserPermissionsListUsers
	}

	if permissions.ManageUsers {
		p |= GroupUserPermissionsManageUsers
	}

	if permissions.ListLabels {
		p |= GroupUserPermissionsListLabels
	}

	if permissions.ManageLabels {
		p |= GroupUserPermissionsManageLabels
	}

	if permissions.ListReviews {
		p |= GroupUserPermissionsListReviews
	}

	if permissions.ManageReviews {
		p |= GroupUserPermissionsManageReviews
	}
	return
}
func getPermissions(p int) (permissions UserPermissions) {
	if p == 0 {
		return
	}

	for i := 0; i < 8; i++ {
		if p&1 == 1 {
			switch 1 << i {
			case GroupUserPermissionsListMedia:
				permissions.ListMedia = true

			case GroupUserPermissionsManageMedia:
				permissions.ManageMedia = true

			case GroupUserPermissionsListUsers:
				permissions.ListUsers = true

			case GroupUserPermissionsManageUsers:
				permissions.ManageUsers = true

			case GroupUserPermissionsListLabels:
				permissions.ListLabels = true

			case GroupUserPermissionsManageLabels:
				permissions.ManageLabels = true

			case GroupUserPermissionsListReviews:
				permissions.ListReviews = true

			case GroupUserPermissionsManageReviews:
				permissions.ManageReviews = true
			}
		}
		p >>= 1
	}
	return
}

func (*GroupInfo) TableName() string {
	return global.DefaultDatabaseGroupTable
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
		err = DB().Table(&gi).Where("groupname", "=", strings.ToLower(i.(string))).Select()
		if gi.Gid == 0 {
			err = errors.New(global.EGroupNotExisted)
		}
	}
	return
}
func groupUpdate(gid int, data interface{}) error {
	_, err := DB().Table(global.DefaultDatabaseGroupTable).Data(data).Where("gid", "=", gid).Update()
	return err
}
func GroupDelete(gid int) error {
	if gid > 1 {
		users, err := GroupGetUsers(gid)
		if err != nil {
			return err
		}

		for k, _ := range users {
			_ = UserRemoveGroup(k, gid)
		}
		_, err = DB().Table(global.DefaultDatabaseGroupTable).Where("gid", "=", gid).Delete()
		return err
	} else {
		return errors.New(global.EGroupForbidden)
	}
}

func GroupGetAll() (gl []GroupInfo, err error) {
	err = DB().Table(&gl).OrderBy("gid").Select()
	return
}

func GroupUpdateMemo(gid int, memo string) error {
	data := map[string]interface{}{"memo": memo}
	return groupUpdate(gid, data)
}
func GroupUpdateName(gid int, gn string) error {
	data := map[string]interface{}{"groupname": gn}
	return groupUpdate(gid, data)
}

func GroupGetMedia(gid int) (mids []int, err error) {
	gl, err := GroupGet(gid)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(gl.Media), &mids)
	return
}
func GroupGetUsers(gid int) (users map[int]int, err error) {
	gi, err := GroupGet(gid)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(gi.Users), &users)
	return
}

func GroupAddMedia(gid, mid int) (err error) {
	mids, err := GroupGetMedia(gid)
	if err != nil {
		return
	}

	for _, v := range mids {
		if mid == v {
			return errors.New(global.EGroupMediaAlreadyInThisGroup)
		}
	}

	mids = append(mids, mid)
	return groupUpdateMedia(gid, mids)
}
func GroupAddUser(gid, uid int, permissions UserPermissions) error {
	fmt.Println("group add user,G=", gid, "U=", uid, "R=", permissions)
	users, err := GroupGetUsers(gid)
	if err != nil {
		return err
	}

	//_, ok := users[uid]
	//if ok {
	//	fmt.Println(global.EGroupUserAlreadyExisted)
	//	return errors.New(global.EGroupUserAlreadyExisted)
	//}

	users[uid] = setPermissions(permissions)
	fmt.Println("user add g")
	err = UserAddGroup(uid, gid)
	if err != nil {
		fmt.Println("E", err.Error())
	}
	fmt.Println("group add u")
	err = groupUpdateUsers(gid, users)
	if err != nil {
		fmt.Println("E", err.Error())
	}
	return err
}

func groupUpdateMedia(gid int, mids []int) error {
	sort.Ints(mids)
	j, err := json.Marshal(mids)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"media": j}
	return groupUpdate(gid, data)
}
func groupUpdateUsers(gid int, users map[int]int) error {
	j, err := json.Marshal(users)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"users": j}
	return groupUpdate(gid, data)
}

func GroupRemoveMedia(gid, mid int) error {
	mids, err := GroupGetMedia(gid)
	if err != nil {
		return err
	}

	for i, v := range mids {
		if v == mid {
			mids = append(mids[:i], mids[i+1:]...)
			return groupUpdateMedia(gid, mids)
		}
	}
	return errors.New(global.EGroupMediaNotExistInGroup)
}
func GroupRemoveUser(gid, uid int) error {
	users, err := GroupGetUsers(gid)
	if err != nil {
		return err
	}

	_, ok := users[uid]
	if !ok {
		return errors.New(global.EGroupUserNotExist)
	}

	delete(users, uid)
	_ = UserRemoveGroup(uid, gid)
	return groupUpdateUsers(gid, users)
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
func GroupSetUserPermissions(gid, uid int, permissions UserPermissions) error {
	users, err := GroupGetUsers(gid)
	if err != nil {
		return err
	}

	_, ok := users[uid]
	if !ok {
		return errors.New(global.EGroupUserNotExist)
	}

	users[uid] = setPermissions(permissions)
	return groupUpdateUsers(gid, users)
}
