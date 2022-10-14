/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiLabel.go
 */

package controller

type LabelData struct {
	MediaHash string `json:"media"`
	Data      string `json:"data"`
	Direction string `json:"direction"`
	Admin     string `json:"admin"`
}

type LabelInfoForButton struct {
	TextBackG string `json:"textbackg"`
	TextHover string `json:"texthover"`
	Tips      string `json:"tips"`
}
