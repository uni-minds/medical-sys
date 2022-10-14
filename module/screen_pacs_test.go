package module

import (
	"testing"
)

func TestPacsGetKeyValue(t *testing.T) {
	data := map[string]PacsKeyValuePair{
		"a": {
			Vr:    "CS",
			Value: []interface{}{"1", "2"},
		},
	}

	v, e := pacsGetKeyValue(data, "a", 1)
	if e != nil {
		t.Log(e.Error())
	} else {
		t.Log(v)
	}
}

func TestPacsImportStudiesAll(t *testing.T) {
	n, err := PacsImportStudiesAll(50)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log("import count:", n)
	}
}

func TestPacsSplitStudiesToGroup(t *testing.T) {
	srcGid := 50
	destGids := []int{51, 52, 53, 54}
	if err := PacsSplitStudiesToGroup(srcGid, destGids); err != nil {
		t.Log(err.Error())
	} else {
		t.Log("OK")
	}
}

func TestPacsGetInstanceImage(t *testing.T) {
	gotBsRaw, gotBsThumb, err := PacsGetInstanceImage("1.2.276.0.26.1.1.1.2.2020.367.3307.4912926.778240", false)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(len(gotBsRaw))
		t.Log(len(gotBsThumb))
	}
}
