package updater

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session
var database *mgo.Database

func ConnectDB() {
	var err error
	session, err = mgo.Dial("192.168.1.6:27017")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connect success.")
	}

	session.SetMode(mgo.Monotonic, true)
	database = session.DB("labelsys")
	return
}

func GetMgo() *mgo.Session {
	return session
}

func GetDB() *mgo.Database {
	return database
}

func DisconnectDB() {
	session.Close()
}

func GetErrNotFound() error {
	return mgo.ErrNotFound
}
