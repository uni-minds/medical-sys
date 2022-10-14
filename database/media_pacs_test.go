package database

import (
	"testing"
)

func TestPacsGetAllStudiesIds(t *testing.T) {
	infos, err := PacsGetAllStudiesIds()
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(infos)
	}
}

func TestPaceInstanceUpdate(t *testing.T) {
	instanceId := "1.2.276.0.26.1.1.1.2.2020.367.3307.4912926.778240"

	info, err := PacsInstanceGetInfo(instanceId)
	if err != nil {
		t.Log("E", err.Error())
	} else {
		t.Log(info.PathCache)
	}

	info.PathCache = "1"
	if err = PacsInstanceUpdate(info); err != nil {
		t.Log("E", err.Error())
	} else {
		t.Log("OK")
	}

	info, _ = PacsInstanceGetInfo(instanceId)
	t.Log(info.PathCache)
}

func TestPacsInstanceUpdateLabel(t *testing.T) {
	instanceId := "1.2.276.0.26.1.1.1.2.2020.367.3307.4912926.778240"
	err := PacsInstanceUpdateLabel(instanceId, "a", "b", "c")
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log("OK")
		t.Log(PacsInstanceGetInfo(instanceId))
	}
}

func TestPacsInstanceGetInfo(t *testing.T) {
	gotInfo, err := PacsInstanceGetInfo("1.2.276.0.26.1.1.1.2.2020.367.2885.4619124.146669568")
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(gotInfo.Frames)
	}

}

func TestPacsSeriesGetInfo(t *testing.T) {
	seriesId := "1.2.276.0.26.1.1.1.2.2020.367.24019.3547836"
	PacsSeriesUpdateAuthor(seriesId, 19, 29)
	i, err := PacsSeriesGetInfo(seriesId)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(i.LabelAuthorUid)
		t.Log(i.LabelProgress)
	}
}
