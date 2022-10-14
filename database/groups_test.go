package database

import (
	"testing"
)

func TestGroupGetAll(t *testing.T) {
	t.Log(GroupGetContains(55))

}

func TestGroupListByUser(t *testing.T) {
	gis, err := GroupListByUser(2)
	t.Log(gis, err)
}

func TestGroupAddUser(t *testing.T) {
	err := GroupAddUser(50, 2, UserPermissions{ListMedia: true})
	t.Log(err)
}
