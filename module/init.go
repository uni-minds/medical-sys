/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: init.go
 */

package module

import (
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
)

var log *logger.Logger

func Init() (err error) {
	log = logger.NewLogger("MODU")
	log.Println("init: module")

	checkDefaultUser()
	return nil
}

func checkDefaultUser() {
	g, err := database.GroupGet(global.DefaultAdminGroup)
	if err != nil {
		gid, _ := database.GroupCreate(database.DbStructGroup{
			Name: global.DefaultAdminGroup,
			Memo: "管理员组",
		})
		log.Println("i", "创建默认管理员组", gid)
		g.Id = gid
	}

	_, err = database.UserGet(global.DefaultAdminUsername)
	if err != nil {
		uid, _ := database.UserCreate(database.DbStructUser{
			Username: global.DefaultAdminUsername,
			Realname: "系统管理员",
			Activate: 1,
			Memo:     "默认管理员",
		})

		p := database.UserPermissions{
			ListUsers:   true,
			ManageUsers: true,
		}
		_ = database.GroupAddUser(g.Id, uid, p)
		_ = UserSetPassword(uid, global.DefaultAdminPassword)
		log.Println("创建默认管理员账户")
	}
}
