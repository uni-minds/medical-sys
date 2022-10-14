/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpUiImport.go
 */

/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  uiImport
 * @Version: 1.0.0
 * @Date: 2020/4/7 12:33
 */

package controller

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/module"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

// /ui/import?path=
// import folder
// -- data.json
// -- files/*.ogv

type JsonData struct {
	DataType string      `json:"type"`
	Data     interface{} `json:"data"`
}

func UiImportMedia(ctx *gin.Context) {
	var dataMedia []module.MediaImportJson
	var dataDicom DicomData
	var bs []byte
	mediaType := "us"
	var jsonData JsonData

	_, uid := CookieValidUid(ctx)

	srcFolder := ctx.Query("path")
	if srcFolder == "" {
		srcFolder = "/data/import"
	}
	srcFolder, err := filepath.Abs(srcFolder)
	if err != nil {
		log("e", "import:", err.Error())
		return
	}

	destFolder, err := filepath.Abs(path.Join(global.GetAppSettings().PathMedia, mediaType, time.Now().Format("20060102-15H")))
	if err != nil {
		log("e", "import:", err.Error())
		return
	}

	log("i", fmt.Sprintf("import from %s => %s", srcFolder, destFolder))

	dirname := filepath.Base(srcFolder)
	targetJson := filepath.Join(srcFolder, "data.json")
	if filepath.Ext(dirname) == ".json" {
		targetJson = srcFolder
	}

	if fp, err := os.Open(targetJson); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return

	} else {
		bs, _ = ioutil.ReadAll(fp)
		fp.Close()
	}

	if err := json.Unmarshal(bs, &jsonData); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return

	} else if jsonData.DataType == "us_media" {
		jbs, _ := json.Marshal(jsonData.Data)
		if err = json.Unmarshal(jbs, &dataMedia); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		} else {
			os.MkdirAll(destFolder, 0777)
			if err := module.MediaImportFromJson(uid, srcFolder, destFolder, dataMedia); err != nil {
				ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))

			} else {
				ctx.JSON(http.StatusOK, SuccessReturn("Import finish"))
			}
		}

	} else if jsonData.DataType == "us_dicom" {
		jbs, _ := json.Marshal(jsonData.Data)
		if err = json.Unmarshal(jbs, &dataDicom); err != nil {
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		} else {
			log("i", "-> group:", dataDicom.Group)
			log("i", "-> keys :", dataDicom.Keywords)
			for studiesId, studiesData := range dataDicom.DicomTree {
				log("i", "+> studies_id:", studiesId)
				for seriesId, seriesData := range studiesData {
					log("i", "-+> series_id:", seriesId)
					for _, instanceId := range seriesData {
						log("i", "--+> insta_id:", instanceId)
					}
				}
			}
		}

	} else {
		ctx.JSON(http.StatusOK, FailReturn(400, "internal error"))
	}

	return
}

type DicomData struct {
	DicomTree map[string]map[string][]string `json:"dicom_tree"`
	Group     string                         `json:"group"`
	Keywords  []string                       `json:"keywords"`
}
