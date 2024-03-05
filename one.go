package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	md52file map[string]string
)

func rmFile(filePath string) {
	os.Remove(filePath)
}

func getFileMd5(filePath string) (string, bool) {
	pFile, err := os.Open(filePath)
	if err != nil {
		return "", false
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil)), true
}

func calMD5(path string) {
	if strings.HasPrefix(path, ".") == true {
		return
	}
	dirList, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range dirList {
		name := item.Name()
		path1 := fmt.Sprintf("%s\\%s", path, name)
		if item.IsDir() {
			if strings.HasPrefix(name, ".") == true {
				continue
			} else {
				calMD5(path1)
			}
		} else {
			md5, ok := getFileMd5(path1)
			if ok {
				path0, ok2 := md52file[md5]
				if ok2 == false {
					md52file[md5] = path1
				} else {
					if len(path0) > len(path1) {
						md52file[md5] = path1
						rmFile(path0)
					} else {
						rmFile(path1)
					}
					log.Printf("%s 和 %s 内容重复\n", path1, path0)
				}
			}
		}
	}
}

func hdlOne(path string) {
	md52file = make(map[string]string, 65536)
	calMD5(path)
}
