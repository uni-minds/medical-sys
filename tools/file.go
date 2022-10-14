/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: file.go
 */

package tools

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
)

func CopyFile(src, dst string) error {
	log("t", "copy:", src, "=>", dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	bar := progressbar.DefaultBytes(sourceFileStat.Size(), "copying")
	io.Copy(io.MultiWriter(destination, bar), source)

	log("t", "finish copy", src)
	return nil
}

func MoveFile(src, dst string) error {
	log("move:", src, "=>", dst)

	if sourceFileStat, err := os.Stat(src); err != nil {
		return err

	} else if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)

	}

	if err := os.Rename(src, dst); err != nil {
		log("e", "move file:", err.Error())
		log("i", "try copy and delete")
		if err = CopyFile(src, dst); err != nil {
			log("e", "copy file:", err.Error())
			return err
		} else if err := os.Remove(src); err != nil {
			log("e", "delete file:", err.Error())
			return err
		}
		log("i", "src deleted")
	}
	log("t", "finish move", src)
	return nil
}
