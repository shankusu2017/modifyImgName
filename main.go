package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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
				go hdlHead(oPath)
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

func doTime(path, name string) {
	nName, ok := shotTimeJPG(path, name)
	if ok {
		hdlRename(path, name, nName)
		return
	}

	nName, ok = earlyTime(path, name)
	if ok {
		hdlRename(path, name, nName)
		return
	}
}

func hdlTime(path string) {
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
				hdlTime(oPath)
			}
		} else {
			doTime(path, name)
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
	pwd = fmt.Sprintf("X:\\img")
	hdlOne(pwd)
	hdlHead(pwd)
	hdlTime(pwd)
	log.Printf("cntRname: %d, cntRemove: %d\n", cntRname, cntRemove)
}
