/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: defaultMenu.go
 */

package global

type MenuStruct struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Controller string       `json:"controller"`
	Icon       string       `json:"icon"`
	Child      []MenuStruct `json:"child"`
}

func DefaultMenuData() []MenuStruct {
	menudata := append(make([]MenuStruct, 0), MenuStruct{
		Id:         "index",
		Name:       "首页",
		Controller: "/ui/home",
		Icon:       "fas fa-th",
	})
	menudata = append(menudata, MenuStruct{
		Name:       "影像筛查",
		Controller: "#",
		Icon:       "fas fa-search",
		Child: []MenuStruct{
			{
				Id:         "screen-us",
				Name:       "超声挑图",
				Icon:       "fas fa-image",
				Controller: "/ui/screen?type=us",
			}, {
				Id:         "screen-ct",
				Name:       "CT",
				Controller: "/ui/screen?type=ct",
				Icon:       "far fa-images",
			}, {
				Id:         "screen-stream",
				Name:       "流媒体",
				Icon:       "fas fa-video",
				Controller: "/ui/screen?type=stream",
			},
		},
	})
	menudata = append(menudata, MenuStruct{
		Name:       "标注检索",
		Controller: "#",
		Icon:       "fas fa-edit",
		Child: []MenuStruct{
			{
				Id:         "label-us",
				Name:       "超声标注",
				Controller: "/ui/medialist?type=us",
				Icon:       "fas fa-notes-medical",
			}, {
				Id:         "label-ct",
				Name:       "CT标注",
				Controller: "/ui/medialist?type=ct",
				Icon:       "fas fa-notes-medical",
			}, {
				Id:         "label-stream",
				Name:       "流媒体",
				Controller: "/ui/medialist?type=steam",
				Icon:       "fas fa-notes-medical",
			},
		},
	})
	//menudata = append(menudata, MenuStruct{
	//	Name:       "分析报告",
	//	Controller: "#",
	//	Icon:       "fas fa-chart-bar",
	//	Child: []MenuStruct{
	//		{
	//			Id:         "analysis-cta",
	//			Name:       "CTA 分析报告",
	//			Controller: "/ui/analysis?type=cta",
	//			Icon:       "fa fa-book",
	//		}, {
	//			Id:         "analysis-ccta",
	//			Name:       "CCTA 分析报告",
	//			Controller: "/ui/analysis?type=ccta",
	//			Icon:       "fa fa-book",
	//		}, {
	//			Id:         "analysis-deepsearch",
	//			Name:       "深度检索报告",
	//			Controller: "/ui/analysis?type=deepsearch",
	//			Icon:       "fa fa-book",
	//		},
	//	},
	//})
	menudata = append(menudata, MenuStruct{
		Name:       "后台管理",
		Controller: "#",
		Icon:       "fas fa-th",
		Child: []MenuStruct{
			{
				Id:         "manage-blockchain",
				Name:       "区块链监控",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/blockchain",
			},
			{
				Id:         "manage-browser",
				Name:       "区块链浏览器",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/browser",
			},
			{
				Id:         "manage-tps-stress",
				Name:       "区块链压力测试",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/blockchain?action=tps-stress",
			},
			{
				Id:         "manage-algo",
				Name:       "算法管理",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/algo",
			},
			{
				Id:         "manage-user",
				Name:       "用户管理",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/user",
			},
			{
				Id:         "management-upload-dicom",
				Name:       "上传-DICOM",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/upload?type=dicom",
			},
		},
	})
	menudata = append(menudata, MenuStruct{
		Id:         "algo-exec",
		Name:       "算法投放",
		Controller: "/mobi",
		Icon:       "fab fa-html5",
	})
	return menudata
}
