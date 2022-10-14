/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: ffmpeg.go
 */

package tools

import (
	"bufio"
	"errors"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os/exec"
	"strconv"
)

func FfConv(input, output, tp string) error {
	kwargs := make(map[string]interface{})
	switch tp {
	case "h264":
		kwargs["b:v"] = "50000K"
		kwargs["vcodec"] = "h264"

	case "gif":
		kwargs["s"] = "640x480"
		kwargs["vcodec"] = "gif"
		kwargs["r"] = "15"

	case "ogv":
		kwargs["b:v"] = "50000K"
		kwargs["vcodec"] = "libtheora"

	}
	return ffmpeg.Input(input).Output(output, kwargs).OverWriteOutput().Run()
}

func FfProbe(input string) (data string, err error) {
	ok, data, _, ec := execCommand("ffprobe", []string{"-show_streams", input})
	if !ok {
		fmt.Println("E")
		return "", errors.New("E")
	} else if ec != 0 {
		return "", errors.New("error code:" + strconv.Itoa(ec))
	}
	return
}

func execCommand(commandName string, params []string) (ok bool, output string, outerr string, exitcode int) {
	cmd := exec.Command(commandName, params...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		ok = false
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		ok = false
		return
	}

	cmd.Start()

	stderrReader := bufio.NewReader(stderr)
	stdoutReader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容

	go func() {
		for {
			line, err2 := stderrReader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			//fmt.Printf(line)
			outerr += line
		}
	}()

	go func() {
		for {
			line, err2 := stdoutReader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			//fmt.Printf(line)
			output += line
		}
	}()

	cmd.Wait()
	exitcode = cmd.ProcessState.ExitCode()
	return true, output, outerr, exitcode
}
