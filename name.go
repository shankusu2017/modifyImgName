package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	wgName sync.WaitGroup
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
			go hdlFile(path, fileName)
		}
	}
}

func hdlFile(path, name string) {
	defer wgName.Done()
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
	nName, ok = shotTimeJPG(path, name)
	if ok == true {
		hdlRename(path, name, nName)
		return
	}
	nName, ok = earlyTime(path, name)
	if ok == true {
		hdlRename(path, name, nName)
		return
	}
}

func hdlName(path string) {
	hdlDir(path)
	wgName.Wait()
}