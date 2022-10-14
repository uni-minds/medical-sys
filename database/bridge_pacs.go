package database

import (
	"gitee.com/uni-minds/bridge_pacs"
	"gitee.com/uni-minds/medical-sys/global"
)

var pacsPort bridge_pacs.DbManager
var pacsServer bridge_pacs.PacsServer

func BridgePacsInit() (err error) {
	app := global.GetAppSettings()
	pacsPort.Init(app.DbFilePacs)
	serverName := "pacs_1"
	log("i", "use pacs db:", app.DbFilePacs)
	if pacsServer, err = pacsPort.GetServer(serverName); err != nil {
		panic(err.Error())
	}
	return err
}

func BridgeGetPacsDatabaseHandler() bridge_pacs.DbManager {
	return pacsPort
}

func BridgeGetPacsServerHandler() bridge_pacs.PacsServer {
	return pacsServer
}
