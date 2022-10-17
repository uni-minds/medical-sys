package module

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"time"
)

func StreamSyncFolder(folder, tagView string) (importd []string, err error) {
	if folder == "" {
		return nil, fmt.Errorf("folder must given")
	}

	importd = make([]string, 0)
	err = filepath.Walk(folder, func(srcFile string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".m3u8" {
			srcFolder := path.Dir(srcFile)
			fmt.Println("import", srcFolder)
			patientId := time.Now().Format("patient_20060102-150405")
			err = MediaImportM3U8(-1, 1, path.Base(srcFolder), srcFolder, "sync_local", tagView, patientId)
			if err == nil {
				importd = append(importd, srcFolder)
			}
			return err
		} else {
			return nil
		}
	})
	return importd, err
}
