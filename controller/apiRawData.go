package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

var menudata []MenuStruct

type MenuStruct struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Controller string       `json:"controller"`
	Icon       string       `json:"icon"`
	Child      []MenuStruct `json:"child"`
}

type ViewStruct struct {
	Default      LabelTool
	Common       []LabelTool
	Quility      []LabelTool
	Time         []LabelTool
	Pathological []LabelTool
}

type LabelTool struct {
	Id     string `json:"id"`     //标注id
	Name   string `json:"name"`   //中文名称
	Group  string `json:"group"`  //组名称
	Color  string `json:"color"`  //主体颜色
	Type   string `json:"type"`   //工具类型
	Domain string `json:"domain"` //作用域 frame / global
	Value  string `json:"value"`  //标签名
	GRadio bool   `json:"gradio"` //组内单选
	GOpen  bool   `json:"gopen"`  //默认展开
}

func GetRawData(ctx *gin.Context) {
	switch ctx.Query("action") {
	// raw?action=getversion
	case "getversion":
		ctx.Writer.WriteString(global.GetCopyrightHtml())

		// raw?action=getmenujson
	case "getmenujson":
		ctx.JSON(http.StatusOK, SuccessReturn(menudata))
		return

		// raw?action=getviewjson&view=4ap
	case "getviewjson":
		view := ctx.Query("view")
		if view != "" {
			ctx.JSON(http.StatusOK, SuccessReturn(getViewData(view)))
		} else {
			ctx.JSON(http.StatusOK, FailReturn("N/A"))
		}
		return

	default:

	}
}

func getViewData(view string) (d []LabelTool) {
	view = strings.ToLower(view)
	d = make([]LabelTool, 0)
	d = []LabelTool{
		{Type: "group", Id: "t", Name: "时间标签", Group: "t", Color: "palevioletred", GRadio: true, GOpen: true},
		{Type: "group", Id: "c", Name: "通用标签", Group: "c", Color: "palevioletred", GRadio: false, GOpen: false},
		{Type: "group", Id: "s", Name: "异常标签", Group: "s", Color: "palevioletred", GRadio: false, GOpen: false},
		{Type: "group", Id: "q", Name: "质量标签", Group: "q", Color: "palevioletred", GRadio: true, GOpen: true},

		{Type: "radio", Group: "t", Domain: "frame", Id: "SSMQ", Name: "收缩末期", Value: "SSMQ"},
		{Type: "radio", Group: "t", Domain: "frame", Id: "SZMQ", Name: "舒张末期", Value: "SZMQ"},
		{Type: "radio", Group: "t", Domain: "frame", Id: "SPEC", Name: "特殊时间", Value: "INPUT"},
		{Type: "radio", Group: "q", Domain: "global", Value: "5", Id: "FQ5", Name: "优秀"},
		{Type: "radio", Group: "q", Domain: "global", Value: "4", Id: "FQ4", Name: "良好"},
		{Type: "radio", Group: "q", Domain: "global", Value: "3", Id: "FQ3", Name: "一般"},
		{Type: "radio", Group: "q", Domain: "global", Value: "2", Id: "FQ2", Name: "差"},
		{Type: "radio", Group: "q", Domain: "global", Value: "1", Id: "FQ1", Name: "不可评估"},
	}

	switch strings.ToLower(view) {
	case "4ap":
		d = append(d, []LabelTool{
			defLabel("XG"), defLabel("JZ"), defLabel("DAO"),
			defLabel("LA"), defLabel("LV"), defLabel("RA"),
			defLabel("RV"), defLabel("二尖瓣前叶"), defLabel("二尖瓣后叶"),
			defLabel("SJBGY"), defLabel("SJBQY"), defLabel("真肋骨1"),
			defLabel("真肋骨2"), defLabel("假肋骨1"), defLabel("假肋骨2"),
			defLabel("心肌外膜"), defLabel("原发房间隔"),

			{Type: "com", Group: "c", Color: "#dbff8e", Id: "FJMZ", Name: "肺静脉左"},
			{Type: "com", Group: "c", Color: "#6487ff", Id: "FJM", Name: "肺静脉右"},
			{Type: "com", Group: "c", Color: "#5BFF7F", Id: "fjgJF", Name: "继发房间隔"},
			{Type: "com", Group: "c", Color: "#FFB0F4", Id: "sjg", Name: "室间隔"},
			{Type: "com", Group: "c", Color: "#ffbd63", Id: "RYKBM", Name: "卵圆孔瓣膜"},
			{Type: "com", Group: "c", Color: "#ff228f", Id: "RYKKK", Name: "卵圆孔开口"}}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#FFA", Id: "TSC1", Name: "TSC1"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "TSC2", Name: "TSC2"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "TSC3", Name: "TSC3"},
			{Type: "com", Group: "s", Color: "#F0E", Id: "CS", Name: "冠状静脉窦"},
			{Type: "com", Group: "s", Color: "#a3eb70", Id: "SJBYC", Name: "三尖瓣异常"},
			{Type: "com", Group: "s", Color: "#a3eb00", Id: "EJBEC", Name: "二尖瓣异常"},
			{Type: "com", Group: "s", Color: "#F00", Id: "GTFSB1", Name: "共同房室瓣1"},
			{Type: "com", Group: "s", Color: "#F01", Id: "GTFSB2", Name: "共同房室瓣2"},
			{Type: "com", Group: "s", Color: "#F02", Id: "FJMGTQ", Name: "肺静脉共同腔"},
			{Type: "com", Group: "s", Color: "#F03", Id: "DXS", Name: "单心室"},
			{Type: "com", Group: "s", Color: "#F04", Id: "DXF", Name: "单心房"}}...)

	case "a", "ac":
		d = append(d, []LabelTool{
			defLabel("JZ"), defLabel("DAO"),
			{Type: "com", Group: "c", Color: "#ffe92c", Id: "UV", Name: "脐静脉"},
			{Type: "com", Group: "c", Color: "#9aff25", Id: "ST", Name: "胃泡"},
			{Type: "com", Group: "c", Color: "#ff3e10", Id: "DV", Name: "静脉导管"},
			{Type: "com", Group: "c", Color: "#FD843F", Id: "IVC", Name: "下腔静脉"},
			{Type: "com", Group: "c", Color: "#FB43FD", Id: "LIVER", Name: "肝脏"},
			{Type: "com", Group: "c", Color: "#ffe2c3", Id: "JJM", Name: "奇静脉"},
			{Type: "com", Group: "c", Color: "#f0e533", Id: "DN", Name: "胆囊"}}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#fcb1a3", Id: "CZJM", Name: "垂直静脉"},
			{Type: "com", Group: "s", Color: "#ffe1a7", Id: "YCXG", Name: "异常血管"},
			{Type: "com", Group: "s", Color: "#FFA", Id: "ERR1", Name: "异常结构1"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "ERR2", Name: "异常结构2"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "ERR3", Name: "异常结构3"}}...)

	case "l":
		d = append(d, []LabelTool{
			defLabel("XG"), defLabel("JZ"), defLabel("DAO"),
			defLabel("LA"), defLabel("LV"), defLabel("AO"),
			defLabel("RA"), defLabel("RV"), defLabel("肋骨1"),
			defLabel("肋骨2"), defLabel("二尖瓣前叶"), defLabel("二尖瓣后叶"),
			defLabel("室间隔"), defLabel("三尖瓣前叶"), defLabel("三尖瓣隔叶"),
			{Type: "com", Group: "c", Color: "#ff7096", Id: "ZDMB1", Name: "主动脉瓣1"},
			{Type: "com", Group: "c", Color: "#5ce63e", Id: "ZDMB2", Name: "主动脉瓣2"},
			{Type: "com", Group: "c", Color: "#ffe1a3", Id: "XJWM", Name: "心肌外膜"}}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#fcb1a3", Id: "EC_ZDMB", Name: "异主动脉瓣"},
			{Type: "com", Group: "s", Color: "#ffe1a7", Id: "EC_EJB", Name: "异二尖瓣"},
			{Type: "com", Group: "s", Color: "#FFA", Id: "TSC1", Name: "TSC1"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "TSC2", Name: "TSC2"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "TSC3", Name: "TSC3"},
			{Type: "com", Group: "s", Color: "#5abd63", Id: "PA", Name: "肺动脉"}}...)

	case "r":
		d = append(d, []LabelTool{
			defLabel("XG"), defLabel("JZ"), defLabel("DAO"),
			defLabel("LA"), defLabel("LV"), defLabel("AO"),
			defLabel("RA"), defLabel("RV"), defLabel("DA"),
			defLabel("PAB1"), defLabel("PAB2"), defLabel("PA"),
			defLabel("肋骨1"), defLabel("肋骨2"), defLabel("SJBGY"),
			defLabel("LPA"), defLabel("RPA"), defLabel("SJBQY")}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#fcb1a3", Id: "EC_PAB1", Name: "异肺动脉瓣"}}...)

	case "van":
		d = append(d, []LabelTool{defLabel("XG"), defLabel("JZ"), defLabel("DAO")}...)
		d = append(d, []LabelTool{defLabel("真肋骨1"), defLabel("真肋骨2"), defLabel("假肋骨1")}...)
		d = append(d, []LabelTool{defLabel("假肋骨2"), defLabel("心肌外膜"), defLabel("肺静脉左")}...)
		d = append(d, []LabelTool{defLabel("肺静脉右"), defLabel("原发房间隔"), defLabel("继发房间隔")}...)
		d = append(d, []LabelTool{defLabel("室间隔"), defLabel("LA"), defLabel("LV")}...)
		d = append(d, []LabelTool{defLabel("RA"), defLabel("RV"), defLabel("二尖瓣前叶")}...)
		d = append(d, []LabelTool{defLabel("二尖瓣后叶"), defLabel("三尖瓣隔叶"), defLabel("三尖瓣前叶")}...)
		d = append(d, []LabelTool{defLabel("卵圆孔瓣膜"), defLabel("卵圆孔开口")}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#FFA", Id: "LJ1", Name: "瘤颈"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "LTCJ1", Name: "瘤体长径"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "PCL1", Name: "膨出瘤"},
			{Type: "com", Group: "s", Color: "#F0E", Id: "P1", Name: "拐点1"},
			{Type: "com", Group: "s", Color: "#a3eb70", Id: "P2", Name: "拐点2"}}...)

	case "3vt", "3v":
		d = append(d, []LabelTool{
			defLabel("XG"), defLabel("JZ"), defLabel("DA"),
			defLabel("真肋骨1"), defLabel("真肋骨2"), defLabel("DAO"),
			defLabel("AO"), defLabel("PA"), defLabel("SVC"),
			defLabel("T"), defLabel("右上腔静脉"), defLabel("无名静脉"),
			defLabel("奇静脉"), defLabel("Thymus")}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#FFA", Id: "YCXG", Name: "异常血管"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "CZJM", Name: "垂直静脉"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "LSVC", Name: "左上腔静脉"},
			{Type: "com", Group: "s", Color: "#FFF", Id: "ERDA", Name: "右位动脉导管"},
		}...)
	}

	fmt.Println("send crf table:\n", d)

	return d
}

func defLabel(name string) LabelTool {
	lt := LabelTool{Type: "com", Group: "c"}
	switch name {
	case "胸骨", "XG":
		lt.Color = "#ffe92c"
		lt.Name = "胸骨"
		lt.Id = "XG"
	case "气管", "T":
		lt.Color = "#f0092c"
		lt.Name = "气管"
		lt.Id = "T"
	case "上腔静脉", "SVC":
		lt.Color = "#f2292c"
		lt.Name = "上腔静脉"
		lt.Id = "SVC"
	case "右上腔静脉", "RSVC":
		lt.Color = "#f4492c"
		lt.Name = "右上腔静脉"
		lt.Id = "RSVC"
	case "无名静脉", "IV":
		lt.Color = "#f6692c"
		lt.Name = "无名静脉"
		lt.Id = "IV"
	case "奇静脉", "VA":
		lt.Color = "#f8892c"
		lt.Name = "奇静脉"
		lt.Id = "VA"
	case "脊柱", "JZ":
		lt.Color = "#9aff25"
		lt.Id = "JZ"
		lt.Name = "脊柱"
	case "降主动脉", "DAO":
		lt.Color = "#ff3e10"
		lt.Id = "DAO"
		lt.Name = "降主动脉"
	case "真肋骨1":
		lt.Color = "#FD843F"
		lt.Id = "LgZ1"
		lt.Name = "真肋骨1"
	case "真肋骨2":
		lt.Color = "#E3FD5B"
		lt.Id = "LgZ2"
		lt.Name = "真肋骨2"
	case "假肋骨1":
		lt.Color = "#709DFD"
		lt.Id = "LgJ1"
		lt.Name = "假肋骨1"
	case "假肋骨2":
		lt.Color = "#FB43FD"
		lt.Id = "LgJ2"
		lt.Name = "假肋骨2"
	case "心肌外膜":
		lt.Color = "#ffe1a3"
		lt.Id = "XJWM"
		lt.Name = "心肌外膜"
	case "肺静脉左":
		lt.Color = "#dbff8e"
		lt.Id = "FJMZ"
		lt.Name = "肺静脉左"
	case "肺静脉右":
		lt.Color = "#6487ff"
		lt.Id = "FJM"
		lt.Name = "肺静脉右"
	case "原发房间隔":
		lt.Color = "#1bdaff"
		lt.Id = "SZ"
		lt.Name = "原发房间隔"
	case "继发房间隔":
		lt.Color = "#5BFF7F"
		lt.Id = "fjgJF"
		lt.Name = "继发房间隔"
	case "室间隔":
		lt.Color = "#FFB0F4"
		lt.Id = "sjg"
		lt.Name = "室间隔"
	case "左房", "LA":
		lt.Color = "#76ff3e"
		lt.Id = "LA"
		lt.Name = "左房"
	case "左室", "LV":
		lt.Color = "#b9ff94"
		lt.Id = "LV"
		lt.Name = "左室"
	case "二尖瓣前叶":
		lt.Color = "#a3ebff"
		lt.Id = "EJBQY"
		lt.Name = "二尖瓣前叶"
	case "二尖瓣后叶":
		lt.Color = "#8bb9ff"
		lt.Id = "EJBHY"
		lt.Name = "二尖瓣后叶"
	case "右房", "RA":
		lt.Color = "#fbaeff"
		lt.Id = "RA"
		lt.Name = "右房"
	case "右室", "RV":
		lt.Color = "#ff37d5"
		lt.Id = "RV"
		lt.Name = "右室"
	case "三尖瓣隔叶", "SJBGY":
		lt.Color = "#ff7096"
		lt.Id = "SJBGY"
		lt.Name = "三尖瓣隔叶"
	case "三尖瓣前叶", "SJBQY":
		lt.Color = "#5ce63e"
		lt.Id = "SJBQY"
		lt.Name = "三尖瓣前叶"
	case "卵圆孔瓣膜":
		lt.Color = "#ffbd63"
		lt.Id = "RYKBM"
		lt.Name = "卵圆孔瓣膜"
	case "卵圆孔开口":
		lt.Color = "#ff228f"
		lt.Id = "RYKKK"
		lt.Name = "卵圆孔开口"
	case "主动脉", "AO":
		lt.Color = "#E3FD5B"
		lt.Id = "AO"
		lt.Name = "主动脉"
	case "动脉导管", "DA":
		lt.Color = "#ff7096"
		lt.Id = "DA"
		lt.Name = "动脉导管"
	case "肋骨1", "LG":
		lt.Color = "#FD843F"
		lt.Id = "LG"
		lt.Name = "肋骨1"
	case "肋骨2":
		lt.Color = "#FD84EF"
		lt.Id = "LG2"
		lt.Name = "肋骨2"
	case "肺动脉瓣1", "PAB1":
		lt.Color = "#ff7096"
		lt.Id = "PAB1"
		lt.Name = "肺动脉瓣1"
	case "肺动脉瓣2", "PAB2":
		lt.Color = "#5ce63e"
		lt.Id = "PAB2"
		lt.Name = "肺动脉瓣2"
	case "肺动脉", "PA":
		lt.Color = "#ffbd63"
		lt.Id = "PA"
		lt.Name = "肺动脉"
	case "左肺动脉", "LPA":
		lt.Color = "#ff7096"
		lt.Id = "LPA"
		lt.Name = "左肺动脉"
	case "右肺动脉", "RPA":
		lt.Color = "#5ce63e"
		lt.Id = "RPA"
		lt.Name = "右肺动脉"
	case "胸腺", "Thymus":
		lt.Color = "#5c003e"
		lt.Id = "Thymus"
		lt.Name = "胸腺"
	default:
		panic(fmt.Sprintf("label info not exist %s\n", name))
	}
	return lt
}

func init() {
	var menuconfig = path.Join(global.GetAppSettings().SystemAppPath, "menu.yaml")
	var data []MenuStruct
	if err := tools.LoadYaml(menuconfig, &data); err != nil {
		menudata = genDefaultMenuData()
		tools.SaveYaml(menuconfig, menudata)
	} else {
		menudata = data
	}
}

func genDefaultMenuData() []MenuStruct {
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
				Id:         "cta-label",
				Name:       "CTA 多专家标注",
				Controller: "http://172.16.1.13/v/",
				Icon:       "fas fa-file-medical-alt",
			}, {
				Id:         "us-label",
				Name:       "超声影像标注",
				Controller: "/ui/labeltool?type=us",
				Icon:       "fas fa-notes-medical",
			},
		},
	})
	menudata = append(menudata, MenuStruct{
		Name:       "分析报告",
		Controller: "#",
		Icon:       "fas fa-chart-bar",
		Child: []MenuStruct{
			{
				Id:         "analysis-cta",
				Name:       "CTA 分析报告",
				Controller: "/ui/analysis?type=cta",
				Icon:       "fa fa-book",
			}, {
				Id:         "analysis-ccta",
				Name:       "CCTA 分析报告",
				Controller: "/ui/analysis?type=ccta",
				Icon:       "fa fa-book",
			}, {
				Id:         "analysis-deepsearch",
				Name:       "深度检索报告",
				Controller: "/ui/analysis?type=deepsearch",
				Icon:       "fa fa-book",
			},
		},
	})
	menudata = append(menudata, MenuStruct{
		Name:       "后台管理",
		Controller: "#",
		Icon:       "fas fa-th",
		Child: []MenuStruct{
			{
				Id:         "management-blockchain",
				Name:       "区块链监控",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/blockchain",
			},
			{
				Id:         "management-browser",
				Name:       "区块链浏览器",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/browser",
			},
			{
				Id:         "manage-user",
				Name:       "用户管理",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/user",
			},
			{
				Id:         "manage-group",
				Name:       "群组管理",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/group",
			},
			{
				Id:         "manage-media",
				Name:       "媒体管理",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/media",
			},
			{
				Id:         "management-upload-dicom",
				Name:       "上传-DICOM",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/upload?type=dicom",
			},
			{
				Id:         "management-upload-us",
				Name:       "上传-超声",
				Icon:       "fas fa-layer-group",
				Controller: "/ui/manage/upload?type=us",
			},
		},
	})
	return menudata
}
