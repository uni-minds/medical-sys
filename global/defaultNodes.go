/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: defaultNodes.go
 * Description:
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
