package tools

import "testing"

func TestCalcResize(t *testing.T) {
	t.Log(CalcResize(1920, 1080, 0, 800))
}
