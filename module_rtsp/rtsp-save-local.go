/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: rtsp-save-local.go
 * Date: 2022/09/21 11:14:21
 */

package module_rtsp

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/module"
	"path"
	"strings"
)

func RtspSaveToDatabase(uriPush, srcFolder string) error {
	var machineId, tagView, tagCustom string

	// pusher.Path = /live/m1/
	// /:machineId/:view/:tag
	for i, content := range strings.Split(uriPush, "/") {
		switch i {
		case 0:
			machineId = content
		case 1:
			tagView = content
		case 2:
			tagCustom = content
		default:
			tagCustom = fmt.Sprintf("%s/%s", tagCustom, content)
		}
	}
	return module.MediaImportM3U8(-1, 1, path.Base(srcFolder), srcFolder, machineId, tagView, tagCustom)
}
