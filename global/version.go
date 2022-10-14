/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: version.go
 * Date: 2022/09/04 08:23:04
 */

package global

import (
	"fmt"
	"time"
)

type Version struct {
	GitCommit string
	Version   string
	BuildTime string
}

var (
	verData       Version
	copyrightHtml string
	versionString string
)

func SetVersion(ver string) {
	verData.Version = ver
}

func GetVersion() string {
	return verData.Version
}

func GenVersionString() {
	versionString = fmt.Sprintf("%s_g%s_b%s", verData.Version, verData.GitCommit, verData.BuildTime)
	copyrightHtml = fmt.Sprintf(`%s<div class="float-right d-none d-sm-inline-block"><b>Ver</b> %s</div>`, DefaultCopyright, verData.Version)
}

func GetVersionString() string {
	if versionString == "" {
		GenVersionString()
	}
	return versionString
}

func SetBuildTime(bt string) {
	t, _ := time.Parse("2006-01-02 15:04:05", bt)
	verData.BuildTime = t.Format("060102t150405")
}

func GetBuildTime() string {
	return verData.BuildTime
}

func SetGitCommit(git string) {
	verData.GitCommit = git
}

func GetGitCommit() string {
	return verData.GitCommit
}

func GetCopyrightHtml() string {
	return copyrightHtml
}
