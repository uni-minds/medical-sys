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
	"path"
	"path/filepath"
	"time"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/module"
)

// /ui/import?path=
// import folder
// -- data.json
// -- files/*.ogv

func UiImportMedia(ctx *gin.Context) {
	var data []module.MediaImportJson
	var bs []byte
	mediaType := "us"

	_, uid := CookieValidUid(ctx)

	srcFolder := ctx.Query("path")
	if srcFolder == "" {
		srcFolder = "./import"
	}
	srcFolder, err := filepath.Abs(srcFolder)
	if err != nil {
		log("e", "import:", err.Error())
		return
	}

	destFolder, err := filepath.Abs(path.Join(global.GetAppSettings().SystemMediaPath, mediaType, time.Now().Format("20060102-15H")))
	if err != nil {
		log("e", "import:", err.Error())
		return
	}

	log("i", "Import media", srcFolder, "=>", destFolder)

	if fp, err := os.Open(filepath.Join(srcFolder, "data.json")); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return

	} else {
		bs, _ = ioutil.ReadAll(fp)
		fp.Close()
	}

	if err := json.Unmarshal(bs, &data); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		return

	} else if len(data) == 0 {
		ctx.JSON(http.StatusOK, FailReturn(400, "empty data.json"))
		return

	} else {
		os.MkdirAll(destFolder, 0777)
	}

	if err := module.MediaImportFromJson(uid, srcFolder, destFolder, data); err != nil {
		ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))

	} else {
		ctx.JSON(http.StatusOK, SuccessReturn("Import finish"))

	}
	return
}
