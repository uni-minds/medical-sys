/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: log.go
 */

package logger

import (
	"compress/gzip"
	"fmt"
	"gitee.com/uni-minds/utils/tools"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const tag = "LOG"

var logger *log.Logger
var fp *os.File

func GetOutput() *os.File {
	if fp == nil {
		fp = os.Stdout
	}
	return fp
}

func Init(logFile string, verbose bool) (err error) {
	logger = log.New()
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{"module", "msg"},
		TimestampFormat: "20060102T150405",
		HideKeys:        true,
		NoColors:        true,
	})

	if verbose {
		logger.SetLevel(log.TraceLevel)
	} else {
		logger.SetLevel(log.WarnLevel)
	}

	if logFile == "" {
		WriteCore(tag, "w", "log -> screen", 1)
		logger.SetOutput(os.Stdout)
		return nil

	} else {
		if logFile, err = filepath.Abs(logFile); err != nil {
			return err
		} else if err = tools.EnsureDir(filepath.Dir(logFile)); err != nil {
			return err
		}

		WriteCore(tag, "w", fmt.Sprintf("log -> file: %s", logFile), 1)

		if s, err := os.Stat(logFile); err == nil && s.Size() > 1000*1024 {
			if err = Archive(logFile); err != nil {
				return err
			}
		}
		fp, _ = os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		logger.SetOutput(fp)
		return nil

	}
}

func Exit() {
	Write(tag, "i", "exit")
	fp.Close()
}

func Archive(logFile string) error {
	logDir := filepath.Dir(logFile)
	gzfile := filepath.Join(logDir, time.Now().Format("archive/log-20060102-150405.gz"))

	if err := tools.EnsureDir(path.Dir(gzfile)); err != nil {
		return err
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

		gw.Header.Name = path.Base(logFile)
		if _, err = io.Copy(gw, src); err != nil {
			return err
		}

		return os.Remove(logFile)
	}
}

func Write(module, level, content string) {
	WriteCore(module, level, content, 3)
}

func WriteCore(module, level, content string, skipCaller int) {
	w := logger.WithFields(map[string]interface{}{"module": module})
	_, f, l, ok := runtime.Caller(skipCaller)
	if ok {
		content = fmt.Sprintf("%s:%d: %s", path.Base(f), l, content)
	}

	switch strings.ToLower(level) {
	case "debug", "d":
		w.Debug(content)

	case "info", "i":
		w.Info(content)

	case "warn", "w":
		w.Warn(content)

	case "error", "e":
		w.Error(content)

	case "fatal", "f":
		w.Fatal(content)

	default:
		w.Trace(content)
	}
}
