package database

import (
	"testing"
)

func TestGroupCreate(t *testing.T) {
	t.Log(GroupCreate(GroupInfo{GroupName: "G1"}))
	t.Log(GroupCreate(GroupInfo{GroupName: "G2"}))
	t.Log(GroupCreate(GroupInfo{GroupName: "G3"}))
	t.Log(GroupAddUser(2, 4, UserPermissions{}))
	t.Log(GroupAddUser(3, 3, UserPermissions{}))
	t.Log(GroupAddUser(4, 3, UserPermissions{}))
}

func TestGroupDelete(t *testing.T) {
	t.Log(GroupAddUser(1, 2, UserPermissions{}))
	t.Log(GroupAddUser(1, 3, UserPermissions{}))
	t.Log(GroupAddUser(1, 4, UserPermissions{}))
	t.Log(GroupDelete(1))
}

func TestGroupGetMedia(t *testing.T) {
	t.Log(GroupGetMedia(1))
	t.Log(GroupAddMedia(1, 2))
	t.Log(GroupAddMedia(1, 4))
	t.Log(GroupGetMedia(1))
	t.Log(GroupRemoveMedia(1, 4))
	t.Log(GroupGetMedia(1))
	t.Log(GroupAddMedia(1, 4))
	t.Log(GroupAddMedia(1, 4))
	t.Log(GroupGetMedia(1))
}

func TestGroupGetUsers(t *testing.T) {
	t.Log(GroupGetUsers(1))
	t.Log(GroupAddUser(1, 2, UserPermissions{
		ManageUsers: true,
	}))
	t.Log(GroupAddUser(1, 3, UserPermissions{
		ManageLabels: true,
	}))
	t.Log(GroupAddUser(1, 4, UserPermissions{
		ManageReviews: true,
	}))
	t.Log(GroupGetUsers(1))
}
