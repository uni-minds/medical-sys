package database

import (
	"encoding/json"
	"errors"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func GroupSetFlagToPacsStudyId(gid int) (err error) {
	data := map[string]interface{}{"group_type": GroupTypePacsStudiesId}
	return groupUpdate(gid, data)
}

func GroupAddPacsStudiesId(gid int, newStudyId string) (err error) {
	studiesIds, err := GroupGetPacsStudiesIds(gid)
	if err != nil {
		return err
	}

	for _, studiesId := range studiesIds {
		if studiesId == newStudyId {
			return nil
		}
	}

	studiesIds = append(studiesIds, newStudyId)
	return GroupSetPacsStudiesIds(gid, studiesIds)
}
func GroupGetPacsStudiesIds(gid int) (studiesIds []string, err error) {
	gi, err := GroupGet(gid)

	if gi.GroupType != GroupTypePacsStudiesId {
		return nil, errors.New("this group is not for pacs|studies id")
	} else if studiesIds, err = tools.StringDecompress(gi.MediaStudiesIDs); err != nil {
		return nil, err
	} else if gi.MediaCounts != len(studiesIds) {
		data := map[string]interface{}{"media_counts": len(studiesIds)}
		groupUpdate(gid, data)
	}
	return studiesIds, nil
}

func GroupSetPacsStudiesIds(gid int, studiesIds []string) (err error) {
	counts := len(studiesIds)
	if strStudyIds, err := json.Marshal(studiesIds); err != nil {
		return err
	} else {
		data := map[string]interface{}{"group_type": GroupTypePacsStudiesId, "media_studies_ids": strStudyIds, "media_counts": counts}
		return groupUpdate(gid, data)
	}
}
