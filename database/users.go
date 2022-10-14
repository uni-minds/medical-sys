/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: users.go
 */

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func (*UserInfo) TableName() string {
	return global.DefaultDatabaseUserTable
}

type UserInfo struct {
	Uid            int    `gorose:"uid"`
	Username       string `gorose:"username"`
	Groups         string `gorose:"groups"`
	Password       string `gorose:"password"`
	PasswordSalt   string `gorose:"passwordsalt"`
	Email          string `gorose:"email"`
	Realname       string `gorose:"realname"`
	RegisterTime   string `gorose:"registertime"`
	Activate       int    `gorose:"activate"`
	ExpireTime     string `gorose:"expiretime"`
	LoginCount     int    `gorose:"logincount"`
	LoginTime      string `gorose:"logintime"`
	LoginFailCount int    `gorose:"loginfailcount"`
	LastGroupId    int    `gorose:"lastGroupId"`
	LastPageIndex  int    `gorose:"lastPageIndex"`
	LastToken      string `gorose:"lastToken"`
	Memo           string `gorose:"memo"`
}

type UserPermissions struct {
	ListMedia     bool
	ManageMedia   bool
	ListUsers     bool
	ManageUsers   bool
	ListLabels    bool
	ManageLabels  bool
	ListReviews   bool
	ManageReviews bool
}

func initUserDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"username" TEXT NOT NULL default "",
	"groups" TEXT NOT NULL default "[]",
	"password" TEXT NOT NULL default "",
	"passwordsalt" TEXT NOT NULL default "",
	"email" TEXT NOT NULL default "",
	"realname" TEXT NOT NULL default "",
	"activate" INTEGER NOT NULL default 0,
	"expiretime" TEXT NOT NULL default "",
	"registertime" TEXT NOT NULL default "",
	"logincount" INTEGET NOT NULL default 0,
	"logintime" TEXT NOT NULL default "",
	"loginfailcount" INTEGER NOT NULL default 0,
	"lastGroupId" INTEGER NOT NULL default 2, 
	"lastPageIndex" INTEGER NOT NULL default 1, 
	"lastToken" TEXT NOT NULL default "",
	"memo" TEXT NOT NULL default "")`, global.DefaultDatabaseUserTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Panic(err.Error())
	}
}
func UserCreate(u UserInfo) (uid int, err error) {
	ut, err := UserGet(u.Username)
	if err != nil {
		u.Uid = 0
		u.Username = strings.ToLower(u.Username)
		u.RegisterTime = time.Now().Format(global.TimeFormat)
		_, err = DB().Table(u.TableName()).Data(u).Insert()
		u, _ = UserGet(u.Username)
		return u.Uid, nil
	}
	return ut.Uid, errors.New(global.EUserAlreadyExisted)
}

func UserGet(i interface{}) (u UserInfo, err error) {
	switch i.(type) {
	case string:
		err = DB().Table(&u).Where("username", "=", strings.ToLower(i.(string))).Select()
	case int:
		err = DB().Table(&u).Where("uid", "=", i).Select()
	}
	if err != nil {
		return
	}

	if u.Uid == 0 {
		err = errors.New(global.EUserNotExist)
	}
	return
}
func UserGetManual(title, content string) (u UserInfo, err error) {
	err = DB().Table(&u).Where(title, "=", content).Select()
	return u, err
}
func UserGetAll() (ul []UserInfo, err error) {
	err = DB().Table(&ul).OrderBy("uid").Select()
	return
}

func userUpdate(uid int, data interface{}) error {
	_, err := DB().Table(global.DefaultDatabaseUserTable).Data(data).Where("uid", "=", uid).Update()
	return err
}
func UserUpdateAccountActiveType(uid int, activeType int) error {
	data := map[string]interface{}{"activate": activeType}
	return userUpdate(uid, data)
}
func UserUpdateTryFailureCount(uid, c int) error {
	data := map[string]interface{}{"loginfailcount": c}
	return userUpdate(uid, data)
}
func UserUpdateLoginCount(uid, c int) error {
	data := map[string]interface{}{"logincount": c, "logintime": time.Now().Format(global.TimeFormat)}
	return userUpdate(uid, data)
}
func UserUpdateLoginExpireTime(uid int, t time.Time) error {
	data := map[string]interface{}{"loginExpireTime": t.Format(global.TimeFormat)}
	return userUpdate(uid, data)
}
func UserUpdatePassword(uid int, passwordC string, passwordSalt string) error {
	data := map[string]interface{}{"password": passwordC, "passwordSalt": passwordSalt}
	return userUpdate(uid, data)
}

func UserSetToken(uid int, token string) error {
	data := map[string]interface{}{"lastToken": token}
	return userUpdate(uid, data)
}

func UserTokenCheck(uid int, token string) bool {
	if token == "" {
		return false
	}

	ui, err := UserGet(uid)
	if err != nil {
		fmt.Println("E:", err.Error())
		return false
	}
	return ui.LastToken == token
}

func UserGetGroups(uid int) (gids []int, err error) {
	ui, err := UserGet(uid)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(ui.Groups), &gids)
	return
}
func userUpdateGroups(uid int, gids []int) error {
	j := "[]"
	if len(gids) > 0 {
		sort.Ints(gids)
		b, err := json.Marshal(gids)
		if err != nil {
			return err
		}
		j = string(b)
	}

	data := map[string]interface{}{"groups": j}
	fmt.Println("after", j)
	return userUpdate(uid, data)
}
func UserAddGroup(uid, gid int) error {
	gids, err := UserGetGroups(uid)
	fmt.Println("before", gids)
	if err != nil {
		return err
	}

	gids = append(gids, gid)
	gids = tools.RemoveDuplicateInt(gids)
	return userUpdateGroups(uid, gids)
}

func UserRemoveGroup(uid, gid int) error {
	gids, err := UserGetGroups(uid)
	if err != nil {
		return err
	}

	gids = tools.RemoveElementInt(gids, gid)
	return userUpdateGroups(uid, gids)

}

func UserGetLastStatus(uid int) (lastGroupId, lastPageIndex int) {
	ui, err := UserGet(uid)
	if err != nil {
		log.Println("DB E:", err.Error())
	}
	return ui.LastGroupId, ui.LastPageIndex
}
func UserSetLastStatus(uid, lastGroupId, lastPageIndex int) error {
	data := map[string]interface{}{"lastGroupId": lastGroupId, "lastPageIndex": lastPageIndex}
	return userUpdate(uid, data)
}

func UserDelete(uid int) (err error) {
	if uid <= 1 {
		return errors.New(global.EUserForbidden)
	}
	_, err = DB().Table(global.DefaultDatabaseUserTable).Where("uid", "=", uid).Delete()
	return
}
