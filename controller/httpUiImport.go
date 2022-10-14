/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: httpUiImport.go
 * Date: 2022/08/15 13:34:15
 */

package controller

// /ui/import?path=
// import folder
// -- data.json
// -- files/*.ogv

type JsonData struct {
	DataType string      `json:"type"`
	Data     interface{} `json:"data"`
}

type DicomData struct {
	DicomTree map[string]map[string][]string `json:"dicom_tree"`
	Group     string                         `json:"group"`
	Keywords  []string                       `json:"keywords"`
}
