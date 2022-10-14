module uni-minds.com/liuxy/medical-sys

go 1.13

require (
	github.com/Unknwon/goconfig v0.0.0-20191126170842-860a72fb44fd
	github.com/antonfisher/nested-logrus-formatter v1.3.0
	github.com/fatih/color v1.9.0
	github.com/gin-gonic/gin v1.6.3
	github.com/gohouse/gorose/v2 v2.1.3
	github.com/gorilla/websocket v1.4.2
	github.com/mattn/go-runewidth v0.0.9
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/nsf/termbox-go v0.0.0-20201124104050-ed494de23a00
	github.com/schollz/progressbar/v3 v3.7.2
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace uni-minds.com/liuxy/medical-sys => ../medical-sys
