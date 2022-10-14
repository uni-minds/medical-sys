package module

import "testing"

func TestD4C_GetStudies(t *testing.T) {
	D4C_GetStudies("192.168.3.101:8080", "dcm4chee-arc/ui2")
}
