/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: defaultLabel.go
 */

package global

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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

func DefaultUltrasonicViewData(view string) (d []LabelTool) {
	view = strings.ToLower(view)
	d = []LabelTool{
		{Type: "group", Id: "t", Name: "切面标签", Group: "t", Color: "palevioletred", GRadio: true, GOpen: true},
		{Type: "group", Id: "c", Name: "通用标签", Group: "c", Color: "palevioletred", GRadio: false, GOpen: false},
		{Type: "group", Id: "s", Name: "异常标签", Group: "s", Color: "palevioletred", GRadio: false, GOpen: false},
		{Type: "group", Id: "q", Name: "质量标签", Group: "q", Color: "palevioletred", GRadio: true, GOpen: true},

		{Type: "radio", Group: "q", Domain: "global", Value: "5", Id: "FQ5", Name: "优秀"},
		{Type: "radio", Group: "q", Domain: "global", Value: "4", Id: "FQ4", Name: "良好"},
		{Type: "radio", Group: "q", Domain: "global", Value: "3", Id: "FQ3", Name: "一般"},
		{Type: "radio", Group: "q", Domain: "global", Value: "2", Id: "FQ2", Name: "差"},
		{Type: "radio", Group: "q", Domain: "global", Value: "1", Id: "FQ1", Name: "不可评估"},
	}
	switch strings.ToLower(view) {
	case "3v", "3vt":
		d = append(d, []LabelTool{
			{Type: "radio", Group: "t", Domain: "frame", Id: "3V", Name: "三血管", Value: "3V"},
			{Type: "radio", Group: "t", Domain: "frame", Id: "3VT", Name: "三血管气管", Value: "3VT"},
			{Type: "radio", Group: "t", Domain: "frame", Id: "SPEC", Name: "其他", Value: "INPUT"},
		}...)

	default:
		d = append(d, []LabelTool{
			{Type: "radio", Group: "t", Domain: "frame", Id: "SSMQ", Name: "收缩末期", Value: "SSMQ"},
			{Type: "radio", Group: "t", Domain: "frame", Id: "SZMQ", Name: "舒张末期", Value: "SZMQ"},
			{Type: "radio", Group: "t", Domain: "frame", Id: "SPEC", Name: "特殊时间", Value: "INPUT"},
		}...)
	}

	switch strings.ToLower(view) {
	case "4ap", "aa", "4cv":
		d = append(d, []LabelTool{
			DefaultUltrasonicLabel("XG"), DefaultUltrasonicLabel("JZ"), DefaultUltrasonicLabel("DAO"),
			DefaultUltrasonicLabel("LA"), DefaultUltrasonicLabel("LV"), DefaultUltrasonicLabel("RA"),
			DefaultUltrasonicLabel("RV"), DefaultUltrasonicLabel("二尖瓣前叶"), DefaultUltrasonicLabel("二尖瓣后叶"),
			DefaultUltrasonicLabel("SJBGY"), DefaultUltrasonicLabel("SJBQY"), DefaultUltrasonicLabel("真肋骨1"),
			DefaultUltrasonicLabel("真肋骨2"), DefaultUltrasonicLabel("假肋骨1"), DefaultUltrasonicLabel("假肋骨2"),
			DefaultUltrasonicLabel("心肌外膜"), DefaultUltrasonicLabel("原发房间隔"),

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
			DefaultUltrasonicLabel("JZ"), DefaultUltrasonicLabel("DAO"),
			{Type: "com", Group: "c", Color: "#ffe92c", Id: "UV", Name: "脐静脉"},
			{Type: "com", Group: "c", Color: "#9aff25", Id: "ST", Name: "胃泡"},
			{Type: "com", Group: "c", Color: "#ff3e10", Id: "DV", Name: "静脉导管"},
			{Type: "com", Group: "c", Color: "#FD843F", Id: "IVC", Name: "下腔静脉"},
			{Type: "com", Group: "c", Color: "#FB43FD", Id: "LIVER", Name: "肝脏"},
			{Type: "com", Group: "c", Color: "#ffe2c3", Id: "JJM", Name: "奇静脉"},
			{Type: "com", Group: "c", Color: "#ffa0c3", Id: "FLK", Name: "腹部轮廓"},
			{Type: "com", Group: "c", Color: "#f0e533", Id: "DN", Name: "胆囊"}}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#fcb1a3", Id: "CZJM", Name: "垂直静脉"},
			{Type: "com", Group: "s", Color: "#ffe1a7", Id: "YCXG", Name: "异常血管"},
			{Type: "com", Group: "s", Color: "#FFA", Id: "ERR1", Name: "异常结构1"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "ERR2", Name: "异常结构2"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "ERR3", Name: "异常结构3"}}...)

	case "l", "lvot":
		d = append(d, []LabelTool{
			DefaultUltrasonicLabel("XG"), DefaultUltrasonicLabel("JZ"), DefaultUltrasonicLabel("DAO"),
			DefaultUltrasonicLabel("LA"), DefaultUltrasonicLabel("LV"), DefaultUltrasonicLabel("AO"),
			DefaultUltrasonicLabel("RA"), DefaultUltrasonicLabel("RV"), DefaultUltrasonicLabel("肋骨1"),
			DefaultUltrasonicLabel("肋骨2"), DefaultUltrasonicLabel("二尖瓣前叶"), DefaultUltrasonicLabel("二尖瓣后叶"),
			DefaultUltrasonicLabel("室间隔"), DefaultUltrasonicLabel("三尖瓣前叶"), DefaultUltrasonicLabel("三尖瓣隔叶"),
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

	case "r", "rvot":
		d = append(d, []LabelTool{
			DefaultUltrasonicLabel("XG"), DefaultUltrasonicLabel("JZ"), DefaultUltrasonicLabel("DAO"),
			DefaultUltrasonicLabel("LA"), DefaultUltrasonicLabel("LV"), DefaultUltrasonicLabel("AO"),
			DefaultUltrasonicLabel("RA"), DefaultUltrasonicLabel("RV"), DefaultUltrasonicLabel("DA"),
			DefaultUltrasonicLabel("PAB1"), DefaultUltrasonicLabel("PAB2"), DefaultUltrasonicLabel("PA"),
			DefaultUltrasonicLabel("肋骨1"), DefaultUltrasonicLabel("肋骨2"), DefaultUltrasonicLabel("SJBGY"),
			DefaultUltrasonicLabel("LPA"), DefaultUltrasonicLabel("RPA"), DefaultUltrasonicLabel("SJBQY")}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#fcb1a3", Id: "EC_PAB1", Name: "异肺动脉瓣"}}...)

	case "van":
		d = append(d, []LabelTool{DefaultUltrasonicLabel("XG"), DefaultUltrasonicLabel("JZ"), DefaultUltrasonicLabel("DAO")}...)
		d = append(d, []LabelTool{DefaultUltrasonicLabel("真肋骨1"), DefaultUltrasonicLabel("真肋骨2"), DefaultUltrasonicLabel("假肋骨1")}...)
		d = append(d, []LabelTool{DefaultUltrasonicLabel("假肋骨2"), DefaultUltrasonicLabel("心肌外膜"), DefaultUltrasonicLabel("肺静脉左")}...)
		d = append(d, []LabelTool{DefaultUltrasonicLabel("肺静脉右"), DefaultUltrasonicLabel("原发房间隔"), DefaultUltrasonicLabel("继发房间隔")}...)
		d = append(d, []LabelTool{DefaultUltrasonicLabel("室间隔"), DefaultUltrasonicLabel("LA"), DefaultUltrasonicLabel("LV")}...)
		d = append(d, []LabelTool{DefaultUltrasonicLabel("RA"), DefaultUltrasonicLabel("RV"), DefaultUltrasonicLabel("二尖瓣前叶")}...)
		d = append(d, []LabelTool{DefaultUltrasonicLabel("二尖瓣后叶"), DefaultUltrasonicLabel("三尖瓣隔叶"), DefaultUltrasonicLabel("三尖瓣前叶")}...)
		d = append(d, []LabelTool{DefaultUltrasonicLabel("卵圆孔瓣膜"), DefaultUltrasonicLabel("卵圆孔开口")}...)

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#FFA", Id: "LJ1", Name: "瘤颈"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "LTCJ1", Name: "瘤体长径"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "PCL1", Name: "膨出瘤"},
			{Type: "com", Group: "s", Color: "#F0E", Id: "P1", Name: "拐点1"},
			{Type: "com", Group: "s", Color: "#a3eb70", Id: "P2", Name: "拐点2"}}...)

	case "3vt", "3v":
		d = append(d, []LabelTool{
			DefaultUltrasonicLabel("XG"), DefaultUltrasonicLabel("JZ"), DefaultUltrasonicLabel("DA"),
			DefaultUltrasonicLabel("真肋骨1"), DefaultUltrasonicLabel("真肋骨2"), DefaultUltrasonicLabel("DAO"),
			DefaultUltrasonicLabel("AO"), DefaultUltrasonicLabel("PA"),
			DefaultUltrasonicLabel("T"), DefaultUltrasonicLabel("右上腔静脉"), DefaultUltrasonicLabel("无名静脉"),
			DefaultUltrasonicLabel("奇静脉"), DefaultUltrasonicLabel("Thymus"), DefaultUltrasonicLabel("SVC"),
			DefaultUltrasonicLabel("LPA"), DefaultUltrasonicLabel("RPA")}...)
		//DefaultUltrasonicLabel("SVC"),

		d = append(d, []LabelTool{
			{Type: "com", Group: "s", Color: "#FFA", Id: "YCXG", Name: "异常血管"},
			{Type: "com", Group: "s", Color: "#FFC", Id: "CZJM", Name: "垂直静脉"},
			{Type: "com", Group: "s", Color: "#FFE", Id: "LSVC", Name: "左上腔静脉"},
			{Type: "com", Group: "s", Color: "#FFF", Id: "ERDA", Name: "右位动脉导管"},
		}...)
	}

	//log("t", "crf table:", d)

	return d
}

func DefaultUltrasonicLabel(name string) LabelTool {
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
		lt.Color = "#33CCFF" //"#a3ebff"
		lt.Id = "EJBQY"
		lt.Name = "二尖瓣前叶"
	case "二尖瓣后叶":
		lt.Color = "#CC00CC" //"#4f61b7"
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
		log("e", "label info not exist:", name)
	}
	return lt
}

func LabelWriteCrf(filename, view string) error {
	data := DefaultUltrasonicViewData(view)

	if data != nil {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		w := csv.NewWriter(f)
		var header = []string{"id", "name", "group", "type", "domain", "value", "color", "gopen", "gradio"}
		w.Write(header)
		for _, d := range data {
			line := []string{d.Id, d.Name, d.Group, d.Type, d.Domain, d.Value, strings.ToUpper(d.Color)}
			if d.GOpen {
				line = append(line, "1")
			} else {
				line = append(line, "")
			}
			if d.GRadio {
				line = append(line, "1")
			} else {
				line = append(line, "")
			}
			w.Write(line)
		}
		w.Flush()
	}
	return nil
}

func LabelCrfFromCsv(view string) error {
	csvRoot := path.Join(GetAppSettings().PathApp, "database", "crf")
	err := filepath.Walk(csvRoot, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		filename := filepath.Base(path)
		filename = strings.Split(filename, ".")[0]
		info := strings.Split(filename, "_")
		fmt.Println(info)
		return nil
	})
	return err

}
