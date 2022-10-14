package database

import (
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"uni-minds.com/medical-sys/global"
)

var once sync.Once
var engin *gorose.Engin

func init() {
	var err error

	once.Do(func() {
		engin, err = gorose.Open(&gorose.Config{
			Driver: "sqlite3",
			Dsn:    global.SqlDBFile,
		})
	})
	if err != nil {
		log.Panic(err.Error())
	}

	initUserDB()
	initMediaDB()
	initGroupDB()
	initLabelDB()

	return
}
func DB() gorose.IOrm {
	return engin.NewOrm()
}
