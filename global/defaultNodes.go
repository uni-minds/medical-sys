/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: defaultNodes.go
 * Date: 2022/08/15 13:21:15
 */

package global

import (
	"net/url"
)

func GetEdaLocalUrl() url.URL {
	return url.URL{
		Scheme:   "http",
		Host:     "localhost",
		RawQuery: getEdaPacsAdminAccess(),
	}
}

func getEdaPacsAdminAccess() string {
	v := url.Values{}
	v.Add("key", "admin")
	return v.Encode()
}
