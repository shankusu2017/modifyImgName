package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type nameFormat struct {
	Head string
	Name string
}

var (
	formats [3]nameFormat
)

func delHead(str string) (string, bool) {
	for _, format := range formats {
		ok := strings.HasPrefix(str, format.Head)
		if ok == false {
			continue
		}
		if len(str) < len(format.Head)+len(format.Name) {
			continue
		}

		return str[len(format.Head):], true
	}

	return "", false
}

func del32CharHead(path, str string) (string, bool) {
	//a0aab87020cd2ca6dbe5323bcbc630ee.mp4
	idx := strings.Index(str, ".")
	if idx != 32 {
		return "", false
	}
	_, ok := hex.DecodeString(str[:idx])
	if ok != nil {
		return "", false
	}

	allName := fmt.Sprintf("%s//%s", path, str)
	finfo, err := os.Stat(allName)
	if err != nil {
		return "", false
	}
	// 获取文件原来的访问时间，修改时间
	// windows下代码如下
	winFileAttr := finfo.Sys().(*syscall.Win32FileAttributeData)
	ftime := time.Unix(winFileAttr.CreationTime.Nanoseconds()/1e9, 0)
	nName := fmt.Sprintf("%4d%02d%02d%02d%02d%02d",
		ftime.Year(), ftime.Month(), ftime.Day(),
		ftime.Hour(), ftime.Minute(), ftime.Second())

	return nName + str[32:], true
}
