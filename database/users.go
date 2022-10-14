/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: users.go
 */

package database

import (
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"strings"
	"time"
)

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
		log.Error(err.Error())
	}
}
func UserCreate(u DbStructUser) (uid int, err error) {
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

func UserGet(i interface{}) (u DbStructUser, err error) {
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
func UserGetManual(title, content string) (u DbStructUser, err error) {
	err = DB().Table(&u).Where(title, "=", content).Select()
	return u, err
}
func UserGetAll() (ul []DbStructUser, err error) {
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

func UserGetLastStatus(uid int) (lastGroupId, lastPageIndex int) {
	ui, err := UserGet(uid)
	if err != nil {
		log.Error(err.Error())
	}
	return ui.LastGroupId, ui.LastPageIndex
}
func UserSetLastStatus(uid, lastGroupId, lastPageIndex int) error {
	data := map[string]interface{}{"lastGroupId": lastGroupId, "lastPageIndex": lastPageIndex}
	return userUpdate(uid, data)
}

func UserDelete(uid int) (err error) {
	if uid <= 1 {
		log.Error(fmt.Sprintf("remove uid failed: uid=%d", uid))
		return errors.New(global.EUserForbidden)
	}
	_, err = DB().Table(global.DefaultDatabaseUserTable).Where("uid", "=", uid).Delete()
	return
}
