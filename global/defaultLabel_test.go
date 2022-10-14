package global

import (
	"testing"
)

func TestDefaultUltrasonicLabel(t *testing.T) {
	t.Log(LabelWriteCrf("3v.csv", "3v"))
	t.Log(LabelWriteCrf("van.csv", "van"))
	t.Log(LabelWriteCrf("r.csv", "r"))
	t.Log(LabelWriteCrf("l.csv", "l"))
	t.Log(LabelWriteCrf("a.csv", "a"))
	t.Log(LabelWriteCrf("4cv.csv", "4cv"))
}

func TestLabelCrfFromCsv(t *testing.T) {
	t.Log(LabelCrfFromCsv("a"))
}
