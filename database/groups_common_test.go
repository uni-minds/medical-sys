package database

import "testing"

func TestGroupCreate(t *testing.T) {
	gi := GroupInfo{
		Gid:             0,
		GroupName:       "pacs_screen",
		DisplayName:     "挑图总组",
		GroupType:       "",
		MediaVideoMids:  "",
		Users:           "",
		Memo:            "超声挑图总目录",
		MediaStudiesIDs: "",
	}
	gid, err := GroupCreate(gi)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log("New group:", gid)
	}
}
