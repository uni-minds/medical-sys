/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: checkMediaUUIDType.go
 * Date: 2022/10/11 16:33:11
 */

package global

import "strings"

const DICOM_TYPE_US_ID = "1.2."
const (
	UUID_TYPE_NORMAL = iota
	UUID_TYPE_DICOM  = iota
)

func GetMediaUUIDType(mediaUUID string) uint {
	if strings.Contains(mediaUUID, DICOM_TYPE_US_ID) {
		return UUID_TYPE_DICOM
	} else {
		return UUID_TYPE_NORMAL
	}
}

func IsDicomUUID(mediaUUID string) bool {
	return GetMediaUUIDType(mediaUUID) == UUID_TYPE_DICOM
}
