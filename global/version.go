/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: version.go
 */

package global

import (
	"fmt"
	"strings"
)

var versionString string
var copyrightHtml string

func SetVersionString(s string) {
	versionString = s
	copyrightHtml = fmt.Sprintf(`<strong>Copyright &copy; 2020 <a href="http://uni-minds.com">Uni-Minds</a> /
<a href="http://www.buaa.edu.cn">Beihang University</a> / 
<a href="http://www.anzhen.org">Beijing Anzhen Hospital, CCMU.</a> /
<a href="http://www.uni-ledger.com">Uni-Ledger Co.,Ltd.</a> / 
<a href="http://www.bijouxhealthcare.cn">Bijoux Healthcare Co.,Ltd.</a></strong> All rights reserved.
<div class="float-right d-none d-sm-inline-block"><b>Ver</b> %s</div>`, strings.Split(s, "_")[0])
}

func GetVersionString() string {
	return versionString
}

func GetCopyrightHtml() string {
	return copyrightHtml
}
