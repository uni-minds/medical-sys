/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: log.go
 */

package logger

import (
	"compress/gzip"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const tag = "LOG"

var logger *log.Logger
var fp *os.File

func init() {
	fmt.Println("init: log")
	logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{})
}

func Init(logFile string, verbose bool) {
	if verbose {
		logger.SetLevel(log.TraceLevel)
		logger.SetFormatter(&nested.Formatter{
			FieldsOrder:     []string{"module", "msg"},
			TimestampFormat: "20060102T150405",
			HideKeys:        true,
		})
	} else {
		s, err := os.Stat(logFile)
		if err == nil && s.Size() > 1000*1024 {
			if err := archive(logFile); err != nil {
				Write(tag, "e", err.Error())
			}
		} else {
			os.MkdirAll(filepath.Dir(logFile), 0600)
		}
		fp, _ = os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		logger.SetOutput(fp)
	}

	//Write(tag, "t", "start")
}

func Exit() {
	Write(tag, "i", "exit")
	fp.Close()
}

func archive(logFile string) error {
	gzfile := filepath.Join(filepath.Dir(logFile), time.Now().Format("archive/log-20060102-150405.gz"))

	if err := os.MkdirAll(path.Dir(gzfile), os.ModePerm); err != nil {
		fmt.Println(err.Error())
	}

	if src, err := os.Open(logFile); err != nil {
		return err
	} else if dest, err := os.Create(gzfile); err != nil {
		src.Close()
		return err
	} else {
		defer src.Close()
		defer dest.Close()

		gw := gzip.NewWriter(dest)
		defer gw.Close()

		gw.Header.Name = "medical-box.log"
		if _, err = io.Copy(gw, src); err != nil {
			return err
		}

		return os.Remove(logFile)
	}
}

func Write(module, level, c string) {
	w := logger.WithField("module", module)
	switch strings.ToLower(level) {
	case "debug", "d":
		w.Debug(c)

	case "info", "i":
		w.Info(c)

	case "warn", "w":
		w.Warn(c)

	case "error", "e":
		w.Error(c)

	case "fatal", "f":
		w.Fatal(c)

	default:
		w.Trace(c)
	}
}
