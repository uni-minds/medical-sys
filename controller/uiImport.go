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
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"uni-minds.com/liuxy/medical-sys/module"
)

// /ui/import?path=
func UIImportMedia(ctx *gin.Context) {
	importPath := "/data/media/medical-sys/us"
	valid, uid := CookieValidUid(ctx)
	if !valid {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	srcFolder := ctx.Query("path")
	if srcFolder == "" {
		srcFolder = "/tmp/imports"
	}
	jsonfile := filepath.Join(srcFolder, "data.json")
	fs, err := os.Stat(jsonfile)
	if err != nil || fs.IsDir() {
		ctx.JSON(http.StatusOK, FailReturn(err.Error()))
		return
	}

	fp, _ := os.OpenFile(jsonfile, os.O_RDONLY, 0777)
	bs, _ := ioutil.ReadAll(fp)

	var data []module.MediaImportJson
	err = json.Unmarshal(bs, &data)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(err.Error()))
		return
	}

	if len(data) == 0 {
		ctx.JSON(http.StatusOK, FailReturn("empty data.json"))
		return
	}

	destFolder := filepath.Join(importPath, time.Now().Format("20060102-15H"))
	os.MkdirAll(destFolder, 0777)

	err = module.MediaImportFromJson(uid, srcFolder, destFolder, data)
	if err != nil {
		ctx.JSON(http.StatusOK, FailReturn(err.Error()))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn("OK"))
	}
	return
}
