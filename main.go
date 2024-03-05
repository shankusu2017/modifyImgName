package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func hdlRename(dir, oName, nName string) {
	if oName == nName {
		return
	}
	oPath := fmt.Sprintf("%s\\%s", dir, oName)
	nPath := fmt.Sprintf("%s\\%s", dir, nName)
	os.Rename(oPath, nPath)
}

func hdlHead(path string) {
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
				hdlHead(oPath)
			}
		} else {
			nName, ok := delHead(item.Name())
			if ok {
				hdlRename(path, name, nName)
				continue
			}

			nName, ok = del32CharHead(path, item.Name())
			if ok {
				hdlRename(path, name, nName)
			}
		}
	}
}

func HdlTime(path string) {
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
				HdlTime(oPath)
			}
		} else {
			nName, ok := shotTimeJPG(path, item.Name())
			if ok {
				hdlRename(path, name, nName)
				continue
			}

			nName, ok = earlyTime(path, item.Name())
			if ok {
				hdlRename(path, name, nName)
			}
		}
	}
}

func init() {
	formats[0].Head = "微信图片_"
	formats[0].Name = "20240226153256"
	formats[1].Head = "WeChat_"
	formats[1].Name = "20240226153302"
	formats[2].Head = "export"
	formats[2].Name = "1709625041"
}

func main() {
	pwd := "D:\\img\\unsync\\24"
	pwd, _ = os.Getwd()
	hdlHead(pwd)
	HdlTime(pwd)
}
