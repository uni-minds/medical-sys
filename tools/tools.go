package tools

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"uni-minds.com/liuxy/medical-sys/global"
)

const EHashNotFile = "target is not file"

var BUFFERSIZE = 1000 * 1000

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
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

	buf := make([]byte, BUFFERSIZE)
	var c int64 = 0
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		c += int64(n)

		fmt.Printf("%d...", int64(c)*100/sourceFileStat.Size())

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	fmt.Printf("OK\n")
	return err
}

func MoveFile(src, dst string) error {
	fmt.Println("Move:", src, "=>", dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	return os.Rename(src, dst)
}

func GetFileMD5(file string) string {
	stat, err := os.Stat(file)
	if err != nil || stat.IsDir() {
		return ""
	}

	filesize := stat.Size()

	fp, err := os.Open(file)
	if err != nil {
		return ""
	}

	buf := make([]byte, BUFFERSIZE)
	m := md5.New()
	c := 0

	fmt.Printf("0...")
	for {
		n, err := fp.Read(buf)
		if err != nil && err != io.EOF {
			return ""
		}

		if n == 0 {
			break
		}

		m.Write(buf[:n])

		c += n
		fmt.Printf("%d...", int64(c)*100/filesize)
	}
	fmt.Printf("OK\n")
	return hex.EncodeToString(m.Sum(nil))
}

func GetStringMD5(str string) (h string) {
	m := md5.New()
	m.Write([]byte(str))
	h = hex.EncodeToString(m.Sum(nil))
	return
}

func GenSaltString(c int) (s string) {
	alphabeta := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < c; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(62))
		t := alphabeta[n.Int64()]
		s += string(t)
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

func FFmpegToH264(input, output string) error {
	ok, _, output, ec := execCommand("ffmpeg", []string{"-i", input, "-b:v", "50000K", "-vcodec", "h264", "-y", output})
	if !ok {
		return errors.New(global.EFFMPEGConvert)
	}
	if ec != 0 {
		fmt.Println(output)
		return errors.New("error code:" + strconv.Itoa(ec))
	}
	return nil
}

func FFmpegToGIF(input, output string) error {
	ok, _, _, _ := execCommand("ffmpeg", []string{"-i", input, "-s", "640x480", "-vcodec", "gif", "-r", "15", "-y", output})
	if !ok {
		return errors.New(global.EFFMPEGConvert)
	}
	return nil
}

func FFmpegToOGV(input, output string) error {
	ok, _, _, _ := execCommand("ffmpeg", []string{"-i", input, "-b:v", "50000K", "-vcodec", "libtheora", "-y", output})
	if !ok {
		return errors.New(global.EFFMPEGConvert)
	}
	return nil
}

func FFprobe(input string) (data string, err error) {
	ok, data, _, ec := execCommand("ffprobe", []string{"-show_streams", input})
	if !ok {
		fmt.Println("E")
		return "", errors.New("E")
	} else if ec != 0 {
		return "", errors.New("error code:" + strconv.Itoa(ec))
	}
	return
}

func RemoveDuplicateInt(a []int) []int {
	i := 0
	for j := 1; j < len(a); j++ {
		if a[i] != a[j] {
			i++
			a[i] = a[j]
		}
	}
	return a[:i+1]
}
func RemoveElementInt(a []int, ele int) []int {
	var o []int
	sort.Ints(a)
	for k, v := range a {
		if ele == v {
			o = append(a[:k], a[k+1:]...)
			return o
		}
	}
	return a
}
