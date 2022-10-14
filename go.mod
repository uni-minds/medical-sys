module uni-minds.com/medical-sys

go 1.13

require (
	github.com/Unknwon/goconfig v0.0.0-20191126170842-860a72fb44fd
	github.com/gin-gonic/gin v1.5.0
	github.com/gohouse/gorose/v2 v2.1.3
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	uni-minds.com/liuxy/medical-sys v0.0.0-00010101000000-000000000000
)

replace uni-minds.com/liuxy/medical-sys => ../medical-sys
