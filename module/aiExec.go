/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: aiExec.go
 */

package module

import "time"

type DeepBuild struct {
	Db         string
	Node       string
	StudiesUID string
	SeriesUID  string
	L          DeepBuildPoint
	R          DeepBuildPoint
	T          string
}

type DeepBuildPoint struct {
	Index int `json:"i"`
	X     int `json:"x"`
	Y     int `json:"y"`
}

type DeepSearch struct {
	Db        string
	Node      string
	StudyUid  string
	SeriesUid string
}

func AiDemoCctaAnalysisDummy(data DeepBuild) (aid string, err error) {
	log("w", "deep build/ccta", data)
	time.Sleep(3 * time.Second)
	return "aid1", nil
}

func AiDemoCctaSearchDummy(data DeepSearch) (sid string, err error) {
	log("w", "deep search/ccta", data)
	time.Sleep(1 * time.Second)
	return "sid1", nil
}

func AiDemoCtaAnalysisDummy(data DeepBuild) (aid string, err error) {
	log("w", "deep build/cta", data)
	time.Sleep(3 * time.Second)
	return "aid1", nil
}

func AiDemoCtaSearchDummy(data DeepSearch) (sid string, err error) {
	log("w", "deep search/cta", data)
	time.Sleep(1 * time.Second)

	return "sid1", nil
}
