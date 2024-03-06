package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

var (
	cntRname  = 0
	cntRemove = 0
	cntMtx    sync.RWMutex
)

const signedCharst = "01234567"
const charset = "0123456789abcdef"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	// 限制第一个字符的范围，使产生的十六进制字符在 int 的范围内
	// 可以直接在前面添加 "-" 而不用担心负数溢出的问题(eg: -f3876898)
	sb.WriteByte(charset[rand.Intn(len(signedCharst))])
	for i := 1; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func hdlRename(dir, oName, nName string) {
	cntMtx.Lock()
	defer cntMtx.Unlock()

	if oName == nName {
		return
	}
	oPath := fmt.Sprintf("%s\\%s", dir, oName)
	nPath := fmt.Sprintf("%s\\%s", dir, nName)
	if DEBUG_MODEL == false {
		exist, err := pathExists(nPath)
		if exist {
			idx := strings.Index(nName, ".") //a.b
			if idx == -1 {
				return
			}
			n := nName[:idx]
			t := nName[idx+1:]
			nName = fmt.Sprintf("%s_%s.%s", n, randomString(8), t)
			nPath = fmt.Sprintf("%s\\%s", dir, nName)
		}
		err = os.Rename(oPath, nPath)
		if err != nil {
			log.Printf("rename.err(%s)  %s ----> %s\n", err.Error(), oPath, nPath)
		}
	}
	//log.Printf("rename  %s ----> %s\n", oPath, nPath)

	cntRname++
	if cntRname%100 == 0 {
		log.Printf("cntRname: %d\n", cntRname)
	}
}

func rmFile(path string) {
	cntMtx.Lock()
	defer cntMtx.Unlock()

	if DEBUG_MODEL {
		log.Printf("remove  %s \n", path)
	} else {
		os.Remove(path)
	}
	cntRemove++
	if cntRemove%100 == 0 {
		log.Printf("cntRemove: %d\n", cntRemove)
	}
}
