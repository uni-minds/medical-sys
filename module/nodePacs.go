/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: nodePacs.go
 */

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/tools"
	"net/url"
)

type PacsSearchParams struct {
	Db               string
	Node             string
	Target           string
	PatientID        string
	StudyInstanceUID string
	SeriesUID        string
	StudyDate        string
	Offset           int
	Limit            int
}

type PacsGetWado struct {
	StudyUID  string
	SeriesUID string
	ObjectUID string
}

type DicomItem struct {
	Value []interface{} `json:"Value"`
	Vr    string        `json:"vr"`
}

type NodePacsInfo struct {
	Nodename string
	DB       string
	Addr     url.URL
}

var pacsResponseData []map[string]DicomItem
var pacsDbs []NodePacsInfo

func UpdatePacsNodes() {
	var list []global.NodeInfo
	var list1 global.NodeInfo
	var resp tools.ConnectData
	var err error

	u := global.GetEdaLocalUrl()
	u.Path = "/api/v1/node/list"
	if resp, _, err = tools.HttpGet(u.String()); err != nil {
		log("e", err.Error())
	} else if resp.Code != 200 {
		log("e", resp.Message)
	}
	bs, _ := json.Marshal(resp.Data)
	if err = json.Unmarshal(bs, &list); err != nil {
		log("e", err.Error())
	}

	u.Path = "/api/v1/node/status"
	if resp, _, err = tools.HttpGet(u.String()); err != nil {
		log("e", err.Error())
	} else if resp.Code != 200 {
		log("e", resp.Message)
	}
	bs, _ = json.Marshal(resp.Data)
	if err = json.Unmarshal(bs, &list1); err != nil {
		log("e", err.Error())
	}

	list = append(list, list1)

	u.Path = "/api/v1/pacs"
	for _, node := range list {
		u.Host = fmt.Sprintf("%s:80", node.IP)

		log("i", "load dbs on eda", u.String())
		if resp, _, err := tools.HttpGet(u.String()); err != nil {
			log("e", err.Error())
		} else {
			var dbs []string
			bs, _ := json.Marshal(resp.Data)
			err := json.Unmarshal(bs, &dbs)
			if err != nil {
				log("e", err.Error())
			} else {
				for _, db := range dbs {
					NodePacsSetAddr(node.Name, db, url.URL{
						Scheme:   "http",
						Host:     u.Host,
						Path:     fmt.Sprintf("/api/v1/pacs/%s", db),
						RawQuery: u.RawQuery,
					})
				}
			}
		}
	}
}

func NodePacsGetAddr(node, db string) (u url.URL, err error) {
	log("i", pacsDbs)
	for _, p := range pacsDbs {
		if p.DB == db && p.Nodename == node {
			return p.Addr, nil
		}
	}
	return u, errors.New("not found")
}

func NodePacsSetAddr(nodename, db string, u url.URL) {
	for i, p := range pacsDbs {
		if p.DB == db && p.Nodename == nodename {
			pacsDbs[i].Addr = u
			return
		}
	}
	pdb := NodePacsInfo{
		Nodename: nodename,
		DB:       db,
		Addr:     u,
	}
	pacsDbs = append(pacsDbs, pdb)
	log("i", "pacsInsert<-", pdb)
}

func NodeSearchPacsDb(data PacsSearchParams) (resp []map[string]DicomItem, err error) {
	//var u url.URL
	//switch node {
	//case "all", "":
	//	resp = make([]map[string]DicomItem, 0)
	//	for _, node := range GetPacsSupportDbNodeNames(db) {
	//		tmp, _ := NodeSearchPacsDb(node, db, data)
	//		resp = append(resp, tmp...)
	//	}
	//	return resp, nil
	//
	//default:
	//	u, err = NodePacsGetAddr(node, db)
	//}
	//log("i", "db", node, db, u)
	//
	//_, bs, err := tools.HttpPost(u.String(), data, "json")
	//
	//if err != nil {
	//	log("e", err.Error())
	//	return nil, err
	//} else if err = json.Unmarshal(bs, &pacsResponseData); err != nil {
	//	log("e", err.Error(), string(bs))
	//} else {
	//	for k, _ := range pacsResponseData {
	//		pacsResponseData[k]["nodename"] = DicomItem{
	//			Value: []interface{}{node},
	//			Vr:    "CS",
	//		}
	//	}
	//}
	return pacsResponseData, err
}

func NodePacsGetInstanceWado(node, lib string, data PacsGetWado) (bs []byte, err error) {
	p := url.Values{}
	p.Set("requestType", "WADO")
	p.Set("studyUID", data.StudyUID)
	p.Set("seriesUID", data.SeriesUID)
	p.Set("objectUID", data.ObjectUID)
	p.Set("contentType", "image/jpeg")
	p.Set("frameNumber", "1")

	u, err := NodePacsGetAddr(node, lib)
	if err != nil {
		return nil, err
	}

	u.Path = "/dcm4chee-arc/aets/AS_RECEIVED/wado"
	u.RawQuery = p.Encode()

	log("i", "get->", u.String())
	_, bs, err = tools.HttpGet(u.String())
	if err != nil {
		log("e", err.Error())
		return nil, err
	} else {
		return bs, nil
	}
}

func GetPacsDbs() (dbs []string) {
	for _, i := range pacsDbs {
		var flag = false
		for _, db := range dbs {
			if i.DB == db {
				flag = true
			}
		}
		if !flag {
			dbs = append(dbs, i.DB)
		}
	}
	return dbs
}

func GetPacsSupportDbNodeNames(db string) (nodes []string) {
	for _, i := range pacsDbs {
		if i.DB == db {
			nodes = append(nodes, i.Nodename)
		}
	}
	return nodes
}
