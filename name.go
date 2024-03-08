package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	wgName           sync.WaitGroup
	routinueFreeName chan bool
)

func hdlDir(path string) {
	dirList, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range dirList {
		fileName := item.Name()
		filePathName := fmt.Sprintf("%s\\%s", path, fileName)
		if item.IsDir() {
			if strings.HasPrefix(fileName, ".") == true {
				continue
			} else {
				hdlDir(filePathName)
			}
		} else {
			wgName.Add(1)
			routinueFreeName <- true
			go hdlFile(path, fileName)
		}
	}
}

func hdlFile(path, name string) {
	defer func() {
		wgName.Done()
		<-routinueFreeName
	}()
	if strings.HasSuffix(name, ".exe") {
		return
	}
	// 去掉固定的前缀
	nName, ok := delHead(name)
	if ok {
		hdlRename(path, name, nName)
	} else {
		nName, ok = del32CharHead(path, name)
		if ok {
			hdlRename(path, name, nName)
		}
	}

	// 尝试用时间来命令某个文件的名字
	nTime, ok := huaweiPhoneTime(path, name)
	if ok {
		hdlRename2Time(path, name, nTime)
		return
	}

	// 尝试读取 jpg 格式中的时间
	nTime, ok = shotTimeJPG(path, name)
	if ok {
		log.Printf("shotTimeJPG done\n")
		hdlRename2Time(path, name, nTime)
		return
	} else {
		nTime, ok = showTime2JPG(path, name)
		if ok {
			log.Printf("---->showTime2JPG<---- done[%s]\n", fmt.Sprintf("%s\\%s", path, name))
			hdlRename2Time(path, name, nTime)
			return
		}
	}
	// 尝试读取文件系统中的时间和文件名中的时间（取较早的那个）
	nTime, ok = earlyTime(path, name)
	if ok {
		log.Printf("earlyTime done\n")
		hdlRename2Time(path, name, nTime)
		return
	}
}

func hdlName(path string) {
	routinueFreeName = make(chan bool, ROUTINUSCNT)
	exif.RegisterParsers(mknote.All...)

	hdlDir(path)
	wgName.Wait()
}
