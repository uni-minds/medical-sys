package database

import (
	"testing"
	"uni-minds.com/liuxy/medical-sys/global"
)

func TestLabelCreate(t *testing.T) {
	li := LabelInfo{
		Uid:  1,
		Mid:  1,
		Memo: "1",
		Type: global.LabelTypeAuthor,
	}
	lid, err := LabelCreate(li)
	t.Log(lid, err)
	lid, err = LabelCreate(li)
	t.Log(lid, err)

	li.Uid = 2
	lid, err = LabelCreate(li)
	t.Log(lid, err)

	li.Type = global.LabelTypeFinal
	lid, err = LabelCreate(li)
	t.Log(lid, err)
}

func TestLabelGetAll(t *testing.T) {
	lis, _ := LabelGetAll(0, 0, "")
	t.Log(lis)

	lis, _ = LabelGetAll(1, 0, global.LabelTypeAuthor)
	t.Log(lis)

	lis, _ = LabelGetAll(0, 1, "")
	t.Log(lis)

	lis, err := LabelGetAll(10, 10, "")
	t.Log(err)
	t.Log(lis)

	lis, _ = LabelGetAll(0, 0, global.LabelTypeAuthor)
	t.Log(lis)
}

func TestLabelUpdateMemo(t *testing.T) {
	lis, _ := LabelGetAll(0, 0, "")
	t.Log(LabelUpdateMemo(lis[0].Lid, "A"))
}
