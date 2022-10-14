/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: init.go
 */

package module

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/medical-sys/tools"
)

func Init() {
	fmt.Println("module init: module")
	progressData = map[int]string{
		0: "待领取",
		1: "正在标注",
		2: "待审核",
		3: "正在审核",
		4: "审核退回",
		5: "待重审",
		6: "",
		7: "审核完成",
	}
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
