/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: apiDatabase.go
 */

package controller

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/module"
	"gitee.com/uni-minds/utils/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)

func PacsSearchProxy(ctx *gin.Context) {
	db := ctx.Param("db")
	node := ctx.Param("node")
	studyUid := ctx.Param("studyuid")
	seriesUid := ctx.Param("seriesuid")
	objectUid := ctx.Param("objectuid")
	addr, err := module.NodePacsGetAddr(node, db)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return
	}

	if objectUid != "" {
		addr.Path = fmt.Sprintf("/api/v1/pacs/%s/%s/%s/%s", db, studyUid, seriesUid, objectUid)
		log("i", "wado->", addr.String())
		_, bs, _ := tools.HttpGet(addr.String())
		ctx.Writer.Write(bs)
		return
	}

	if seriesUid != "" {
		addr.Path = fmt.Sprintf("/api/v1/pacs/%s/%s/%s", db, studyUid, seriesUid)
	} else if studyUid != "" {
		addr.Path = fmt.Sprintf("/api/v1/pacs/%s/%s", db, studyUid)
	} else {
		addr.Path = fmt.Sprintf("/api/v1/pacs/%s", db)
	}
	u := addr.String()

	resp, _, err := tools.HttpPost(u, nil, "json")
	//log("i",resp)
	if resp.Code == 200 {
		var data []map[string]interface{}
		bs, _ := json.Marshal(resp.Data)
		json.Unmarshal(bs, &data)
		//log("i", "data->", data)
		value := []string{node}
		nodedata := map[string]interface{}{"Value": value, "vr": "CS"}
		for i, ele := range data {
			ele["nodename"] = nodedata
			data[i] = ele
		}
		ctx.JSON(http.StatusOK, SuccessReturn(data))
	} else {
		ctx.JSON(http.StatusOK, FailReturn(resp.Code, resp.Message))
	}
}

func PacsWado(ctx *gin.Context) {
	node := global.GetEdaLocalUrl()
	p := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = node.Scheme
			req.URL.Host = node.Host
			req.URL.RawPath = "/api/v1/pacs/wado"
			req.Host = node.Host
		},
	}
	p.ServeHTTP(ctx.Writer, ctx.Request)
}

func PacsGet(c *gin.Context) {
	fmt.Println("pacs get")
	module.UpdatePacsNodes()
	dbs := module.GetPacsDbs()
	c.JSON(http.StatusOK, SuccessReturn(dbs))
}

func PacsGetNodelist(ctx *gin.Context) {
	db := ctx.Param("db")
	nodes := module.GetPacsSupportDbNodeNames(db)
	ctx.JSON(http.StatusOK, SuccessReturn(nodes))
}

func PacsSearchNode(ctx *gin.Context) {

	var data module.PacsSearchParams
	if err := ctx.BindJSON(&data); err != nil {
		log("e", err.Error())
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return
	}

	data.Db = ctx.Param("db")
	data.Node = ctx.Param("node")
	data.StudyInstanceUID = ctx.Param("studyuid")
	data.SeriesUID = ctx.Param("seriesuid")

	if resp, err := module.NodeSearchPacsDb(data); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(resp))
	}
	return
}

func DatabaseGetNodeObjectUid(c *gin.Context) {
	db := c.Param("db")
	node := c.Param("node")

	data := module.PacsGetWado{
		SeriesUID: c.Param("seriesuid"),
		StudyUID:  c.Param("studyuid"),
		ObjectUID: c.Param("objectuid"),
	}
	bs, err := module.NodePacsGetInstanceWado(node, db, data)

	if err != nil {
		log("e", "wado:", err.Error())
	}
	c.Writer.Write(bs)
	return
}
