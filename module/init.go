/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: init.go
 */

package module

import (
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/logger"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func Init() {
	checkDefaultUser()
}

func checkDefaultUser() {
	g, err := database.GroupGet(global.DefaultAdminGroup)
	if err != nil {
		gid, _ := database.GroupCreate(database.GroupInfo{
			GroupName: global.DefaultAdminGroup,
			Memo:      "管理员组",
		})
		log("i", "创建默认管理员组", gid)
		g.Gid = gid
	}

	_, err = database.UserGet(global.DefaultAdminUsername)
	if err != nil {
		uid, _ := database.UserCreate(database.UserInfo{
			Username: global.DefaultAdminUsername,
			Realname: "系统管理员",
			Activate: 1,
			Memo:     "默认管理员",
		})

		p := database.UserPermissions{
			ListUsers:   true,
			ManageUsers: true,
		}
		_ = database.GroupAddUser(g.Gid, uid, p)
		_ = UserSetPassword(uid, global.DefaultAdminPassword)
		log("i", "创建默认管理员账户")
	}
}

const tag = "MODU"

func log(level string, message ...interface{}) {
	msg := tools.ExpandInterface(message)
	logger.Write(tag, level, msg)
}
