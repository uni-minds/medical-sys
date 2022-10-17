package module

import "testing"

func TestStreamSyncFolder(t *testing.T) {
	str, err := StreamSyncFolder("/Users/liuxy/go/src/gitee.com/uni-minds/medical-sys/tmp/rtsp", "tag")
	if err != nil {
		t.Errorf(err.Error())
	} else {
		t.Log(str)
	}
}
