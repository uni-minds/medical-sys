/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: checksum.go
 */

package tools

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path"
)

func GetFileMD5(file string) string {
	stat, err := os.Stat(file)
	if err != nil || stat.IsDir() {
		return ""
	}

	fp, err := os.Open(file)
	if err != nil {
		return ""
	}

	m := md5.New()

	bar := progressbar.DefaultBytes(stat.Size(), "checksum")
	_, err = io.Copy(io.MultiWriter(m, bar), fp)
	if err != nil {
		log("e", err.Error())
		return ""
	} else {
		checksum := hex.EncodeToString(m.Sum(nil))
		log("finish md5", checksum, "<=", path.Base(file))
		return checksum
	}
}

func GetStringMD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
