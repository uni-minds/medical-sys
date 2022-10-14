package module

import (
	"uni-minds.com/medical-sys/database"
)

const (
	GroupLevelCustom     = -1
	GroupLevelNotAllowed = 0
	GroupLevelGuest      = 1
	GroupLevelMember     = 2
	GroupLevelLeader     = 3
	GroupLevelMaster     = 4
)

func GroupCreate(groupname, displayname, memo string) error {
	gi := database.GroupInfo{
		GroupName:   groupname,
		DisplayName: displayname,
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
	media, err := database.GroupGetMedia(gid)
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

func GroupAddUser(gid int, uid int, role string) error {
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
	}
	return database.GroupAddUser(gid, uid, permissions)
}
