package global

import "time"

func getDefaultConfig() AppSettings {
	return AppSettings{
		Paths: Paths{
			Application: "/data/medisys/application",
			Media:       "/data/medisys/media",
			Cache:       "/data/medisys/cache",
			Log:         "/data/medisys/log",
			Database:    "/data/medisys/database",
		},
		UserRegister: UserRegister{
			Enable:  true,
			RegCode: "BUAA@2022",
		},
		Ports: Ports{
			HTTP: 80,
			RTSP: 554,
			RPC:  8096,
		},
		CookieMaxAge: 24 * int(time.Hour.Seconds()),
		Rtsp: Rtsp{
			Timeout:              0,
			NetworkBuffer:        204800,
			CloseOld:             false,
			PlayerQueueLimit:     0,
			DropPacketWhenPaused: false,
			GopCacheEnable:       true,
			DebugLogEnable:       false,
			TsDurationSecond:     6,
			UploadRoot:           "./tmp/rtsp",
			SaveStreamToLocal:    true,
			AuthorizationEnable:  false,
			ClientUser:           "admin",
			ClientPassword:       "admin",
			FFmpegEncoder:        "-c:v copy -c:a copy",
		},
	}
}
