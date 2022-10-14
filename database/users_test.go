package database

import (
	"testing"
)

func TestUserCreate(t *testing.T) {
	t.Log(UserCreate(UserInfo{Username: "teSt1"}))
	t.Log(UserCreate(UserInfo{Username: "tEst2"}))
	t.Log(UserCreate(UserInfo{Username: "Test3"}))
}

func TestUserDelete(t *testing.T) {
	t.Log(UserDelete(1))
	t.Log(UserDelete(2))
}

func TestUserUpdateGroups(t *testing.T) {
	userUpdateGroups(1, nil)
	t.Log(UserAddGroup(1, 2))
	t.Log(UserGetGroups(1))
	t.Log(UserAddGroup(1, 3))
	t.Log(UserGetGroups(1))
	t.Log(UserRemoveGroup(1, 3))
	t.Log(UserGetGroups(1))
	t.Log(UserAddGroup(1, 1))
	t.Log(UserGetGroups(1))
}
