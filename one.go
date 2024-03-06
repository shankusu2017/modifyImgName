package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	file2md5times map[string]int // 文件内容对应的md5第N次出现
	md5Times      map[string]int // 文件内容对应的md5出现的次数
	mtxMD5TIMES   sync.RWMutex
	wgMD5         sync.WaitGroup
)

func addMD5Times(file, md5 string) int {
	mtxMD5TIMES.Lock()
	defer mtxMD5TIMES.Unlock()

	val := md5Times[md5]
	val++
	md5Times[md5] = val
	file2md5times[file] = val

	return val
}

func getFileMd5(filePath string) {
	defer wgMD5.Done()
	pFile, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	md5str := hex.EncodeToString(md5h.Sum(nil))
	addMD5Times(filePath, md5str)
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
			wgMD5.Add(1)
			go getFileMd5(path1)
		}
	}
}

func delFile(path string, times int) {
	defer wgMD5.Done()
	if times <= 1 {
		return
	}

	os.Remove(path)
}

func hdlSame() {
	for k, v := range file2md5times {
		wgMD5.Add(1)
		go delFile(k, v)
	}
}

func hdlOne(path string) {
	md5Times = make(map[string]int, 65536)
	file2md5times = make(map[string]int, 65536)
	calMD5(path)
	wgMD5.Wait()
	hdlSame()
	wgMD5.Wait()
	log.Printf("md5 cal done!\n")
}
