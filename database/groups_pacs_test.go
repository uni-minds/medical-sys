package database

import (
	"testing"
)

func TestGroupSetPacsStudyIds(t *testing.T) {
	//studyIds := []string{"a", "b", "c"}
	err := GroupSetPacsStudiesIds(50, nil)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log("OK1")
	}

	//err = GroupAddPacsStudiesId(50,"e")
	//if err != nil {
	//	t.Log(err.Error())
	//} else {
	//	t.Log("OK2")
	//}
}

func TestGroupGetPacsStudyIds(t *testing.T) {
	gotStudyIds, err := GroupGetPacsStudiesIds(50)
	t.Log(gotStudyIds, err)
}
