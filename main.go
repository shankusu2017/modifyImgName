package main

import (
	"encoding/hex"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"log"
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
	formats [2]nameFormat
)

func shotTimeJPG(path, str string) (string, bool) {
	allName := fmt.Sprintf("%s//%s", path, str)

	if strings.HasSuffix(str, ".jpg") == false {
		if strings.HasSuffix(str, ".JPG") == false {
			return "", false
		}
	}

	sufIdx := strings.LastIndex(str, ".")

	f, err := os.Open(allName)
	if err != nil {
		log.Printf(err.Error())
		return "", false
	}
	defer f.Close()
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Printf(err.Error())
		return "", false
	}

	//camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
	//fmt.Println(camModel.StringVal())
	//
	//focal, _ := x.Get(exif.FocalLength)
	//numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	//fmt.Printf("%v/%v", numer, denom)

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, _ := x.DateTime() // 拍摄时间
	if tm.Year() <= 2000 {
		return "", false
	}

	nName := fmt.Sprintf("%4d%02d%02d%02d%02d%02d%s",
		tm.Year(), tm.Month(), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second(), str[sufIdx:])

	return nName, true
}

func earlyTime(path, str string) (string, bool) {
	return "", false
}

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

func hdlRename(dir, oName, nName string) {
	oPath := fmt.Sprintf("%s\\%s", dir, oName)
	nPath := fmt.Sprintf("%s\\%s", dir, nName)
	os.Rename(oPath, nPath)
}

func hdlDir(path string) {
	if strings.HasPrefix(path, ".") == true {
		return
	}
	dirList, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range dirList {
		name := item.Name()
		oPath := fmt.Sprintf("%s\\%s", path, name)
		if item.IsDir() {
			if strings.HasPrefix(name, ".") == true {
				continue
			} else {
				hdlDir(oPath)
			}
		} else {
			nName, ok := shotTimeJPG(path, item.Name())
			if ok {
				hdlRename(path, name, nName)
			} else {
				nName, ok = delHead(item.Name())
				if ok {
					hdlRename(path, name, nName)
				} else {
					nName, ok = del32CharHead(path, item.Name())
					if ok {
						hdlRename(path, name, nName)
					} else {
						nName, ok = earlyTime(path, item.Name())
						if ok {
							hdlRename(path, name, nName)
						}
					}
				}
			}
		}
	}
}

func main() {
	formats[0].Head = "微信图片_"
	formats[0].Name = "20240226"
	formats[1].Head = "WeChat_"
	formats[1].Name = "20240226"

	pwd := "D:\\img\\unsync\\24"
	pwd, _ = os.Getwd()
	hdlDir(pwd)
}
