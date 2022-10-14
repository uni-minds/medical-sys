/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: defaultAlgorithms.go
 */

package global

import "fmt"

type AlgorithmInfo struct {
	Index int    `json:"algo-index"`
	Name  string `json:"algo-name"`
	Ref   string `json:"algo-ref"`
}

func DefaultAlgorithms() (algos []AlgorithmInfo) {
	algoServer := "algo-master"
	algos = []AlgorithmInfo{{
		Index: 1,
		Name:  "算法1-MD5校验",
		Ref:   fmt.Sprintf("%s/md5sum:v1", algoServer),
	}, {
		Index: 2,
		Name:  "算法3-甲状腺结节识别",
		Ref:   fmt.Sprintf("%s/thyroid-nodule-detection:v1", algoServer),
	}, {
		Index: 3,
		Name:  "功能1-图片编码上链（BASE64）",
		Ref:   "#v-algo-pic-encode-base64",
	}}

	return algos
}
