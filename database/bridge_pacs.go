package database

import (
	"gitee.com/uni-minds/bridge-pacs/pacs_server"
	"gitee.com/uni-minds/medical-sys/global"
)

var pacsManager *pacs_server.PacsManager
var pacsServer *pacs_server.PacsServer

func BridgePacsInit() {
	dbfile, err := global.GetDbFile("pacs")
	if err != nil {
		log.Error(err.Error())
	}
	pacsManager, _ = pacs_server.NewPacsDbManager(dbfile, global.FlagGetVerbose())

	serverName := "pacs_1"
	log.Println("pacs db ->", dbfile, serverName)
	if pacsServer, err = pacsManager.GetServer(serverName); err != nil {
		panic(err.Error())
	}
}

func BridgeGetPacsDatabaseHandler() *pacs_server.PacsManager {
	return pacsManager
}

func BridgeGetPacsServerHandler() *pacs_server.PacsServer {
	return pacsServer
}
