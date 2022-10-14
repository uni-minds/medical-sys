/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: defaultMenu.go
 * Description:
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
		Name:       "影像检索",
		Controller: "#",
		Icon:       "fas fa-search",
		Child: []MenuStruct{
			{
				Id:         "ct-medialist",
				Name:       "CT 影像检索",
				Controller: "/ui/medialist?type=ct",
				Icon:       "far fa-images",
			}, {
				Id:         "us-medialist",
				Name:       "超声影像检索",
				Icon:       "fas fa-image",
				Controller: "/ui/medialist?type=us",
			},
		},
	})
	menudata = append(menudata, MenuStruct{
		Name:       "标注系统",
		Controller: "#",
		Icon:       "fas fa-edit",
		Child: []MenuStruct{
			{
				//	Id:         "cta-label",
				//	Name:       "CTA 多专家标注",
				//	Controller: "http://172.16.1.13/v/",
				//	Icon:       "fas fa-file-medical-alt",
				//}, {
				Id:         "us-label",
				Name:       "超声影像标注",
				Controller: "/ui/medialist?type=us",
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
			//{
			//	Id:         "manage-group",
			//	Name:       "群组管理",
			//	Icon:       "fas fa-layer-group",
			//	Controller: "/ui/manage/group",
			//},
			//{
			//	Id:         "manage-media",
			//	Name:       "媒体管理",
			//	Icon:       "fas fa-layer-group",
			//	Controller: "/ui/manage/media",
			//},
			{
				Id:         "management-upload-dicom",
				Name:       "上传-DICOM",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/upload?type=dicom",
			},
			//{
			//	Id:         "management-upload-us",
			//	Name:       "上传-超声",
			//	Icon:       "fas fa-layer-group",
			//	Controller: "/ui/manage/upload?type=us",
			//},
		},
	})
	menudata = append(menudata, MenuStruct{
		Id:         "H5",
		Name:       "H5",
		Controller: "/mobi",
		Icon:       "fab fa-html5",
	})
	return menudata
}
