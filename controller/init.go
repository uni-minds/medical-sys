package controller

import (
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/tools"
	"path"
)

func Init() {
	algofile = path.Join(global.GetAppSettings().PathApp, "algo.yaml")
	if err := tools.LoadYaml(algofile, &algolist); err != nil {
		algolist = global.DefaultAlgorithms()
		tools.SaveYaml(algofile, algolist)
	}

	var menuconfig = path.Join(global.GetAppSettings().PathApp, "menu.yaml")
	if err := tools.LoadYaml(menuconfig, &menudata); err != nil {
		menudata = global.DefaultMenuData()
		tools.SaveYaml(menuconfig, menudata)
	}
}
