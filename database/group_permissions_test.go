package database

import "testing"

func Test_setPermissions(t *testing.T) {
	p := UserPermissions{}
	role := "master"

	p.ListMedia = true

	if role == "member" || role == "leader" || role == "master" {
		p.ListLabels = true
		p.ListUsers = true
	}

	if role == "leader" || role == "master" {
		p.ManageReviews = true
		p.ListReviews = true
	}

	if role == "master" {
		p.ManageLabels = true
		p.ManageReviews = true
		p.ManageUsers = true
	}
	P := setPermissions(p)
	t.Log(P)
}
