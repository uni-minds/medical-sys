module gitee.com/uni-minds/medical-sys

go 1.17

replace gitee.com/uni-minds/bridge_pacs => ../bridge_pacs

replace gitee.com/uni-minds/bridge_his => ../bridge_his

require (
	gitee.com/uni-minds/bridge_his v0.0.0-20210810064955-595eb84f7418
	gitee.com/uni-minds/bridge_pacs v0.0.0-20210818030517-70c5f0f81cb1
	github.com/Unknwon/goconfig v1.0.0
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/fatih/color v1.13.0
	github.com/gin-gonic/gin v1.7.3
	github.com/gohouse/gorose/v2 v2.1.12
	github.com/mattn/go-runewidth v0.0.13
	github.com/mattn/go-sqlite3 v1.14.9
	github.com/nsf/termbox-go v1.1.1
	github.com/schollz/progressbar/v3 v3.8.3
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/gohouse/golib v0.0.0-20210711163732-a5c22059eb75 // indirect
	github.com/gohouse/t v0.0.0-20201007094014-630049a6bfe9 // indirect
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
