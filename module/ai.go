package module

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

type CCTA_PARAMS struct {
}

// 模型->特征向量
func RunAlgo(algo, params string) (aid string, err error) {
	switch algo {
	case "ccta_get_features":
		return algoCctaGetFeatures()

	case "cta_get_features":
		return algoCtaGetFeatures()
	default:
		return "", errors.New("unknown algo")
	}
}

func algoCctaGetFeatures() (aid string, err error) {
	aid = tools.RandString(32)
	log("w", "exec: ccta extra")
	return aid, nil
}

func algoCtaGetFeatures() (aid string, err error) {
	aid = tools.RandString(32)
	log("w", "exec: cta extra")
	return aid, nil
}

// 特征向量->输出

func AlgoCctaGetFeatureResult(aid, part string) (result map[string]string) {
	result = make(map[string]string)

	fileDir := path.Join(global.GetAppSettings().SystemAppPath, "ai_data/ccta/result/json")
	if part != "" {
		file := path.Join(fileDir, fmt.Sprintf("%s.json", strings.ToLower(part)))
		log("i", "get part", part, "from file:", file)
		if fp, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm); err != nil {
			log("e", err.Error())
		} else if bs, err := ioutil.ReadAll(fp); err != nil {
			log("e", err.Error())
		} else {
			result[part] = string(bs)
			return result
		}
		return nil
	} else {
		files := make([]string, 0)
		err := filepath.Walk(fileDir, func(filename string, fi os.FileInfo, err error) error {
			if fi.IsDir() {
				return nil
			} else if strings.HasSuffix(strings.ToLower(fi.Name()), ".json") {
				files = append(files, filename)
			}
			return nil
		})

		if err != nil {
			log("e", err.Error())
			return nil
		}

		for _, file := range files {
			if fp, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm); err != nil {
				log("e", err.Error())
			} else if bs, err := ioutil.ReadAll(fp); err != nil {
				log("e", err.Error())
			} else {
				filename := path.Base(file)
				base := strings.TrimSuffix(filename, path.Ext(filename))
				_, ok := result[base]
				if ok {
					log("e", "result has the same index", base)
				} else {
					result[base] = string(bs)
				}
			}
		}
		return result
	}
}
