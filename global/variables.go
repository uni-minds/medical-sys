package global

import "log"

type AppSettings struct {
	EnableUserRegister bool
}

var mediaRoot string
var cookieMaxAge int
var appSettings AppSettings
var regcode string

func SetMediaRoot(root string) {
	mediaRoot = root
}

func GetMediaRoot() string {
	return mediaRoot
}

func SetCookieMaxAge(second int) {
	log.Println("cookie max age =", second, "seconds.")
	cookieMaxAge = second
}

func GetCookieMaxAge() int {
	return cookieMaxAge
}

func GetAppSettings() AppSettings {
	return appSettings
}

func SetAppSettings(s AppSettings) {
	appSettings = s
}

func SetUserRegCode(s string) {
	regcode = s
}

func GetUserRegCode() string {
	return regcode
}
