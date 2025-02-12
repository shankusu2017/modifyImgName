package main

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type nameFormat struct {
	Head string
	Name string
}

var (
	formats [1024]nameFormat
)

func init() {
	formats[0].Head = "微信图片_"
	formats[0].Name = "20240226153256"
	formats[1].Head = "WeChat_"
	formats[1].Name = "20240226153302"
	formats[2].Head = "export"
	formats[2].Name = "1709625041"
	formats[3].Head = "mmexport"
	formats[3].Name = "1709625041"
	formats[4].Head = "IMG_"
	formats[4].Name = "20250211_191223"
}

func delHead(str string) (string, bool) {
	for _, format := range formats {
		/* 所有的模式均已匹配完毕 */
		if len(format.Head) == 0 {
			break
		}
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
	_, err := hex.DecodeString(str[:idx])
	if err != nil {
		return "", false
	}

	ftime, ok := fileSysTime(path, str)
	if ok == true {
		nName := fmt.Sprintf("%4d%02d%02d%02d%02d%02d",
			ftime.Year(), ftime.Month(), ftime.Day(),
			ftime.Hour(), ftime.Minute(), ftime.Second())
		return nName + str[32:], true
	} else {
		return "", false
	}
}
