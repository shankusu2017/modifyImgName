package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

func isValidData(year, month, day, hour, min, sec int) bool {
	if (year < 2000 || year > 3000) ||
		(month < 1 || month > 12) ||
		(day < 1 || day > 31) ||
		(hour < 0 || hour > 23) ||
		(min < 0 || min > 59) ||
		(sec < 0 || sec > 59) {
		return false
	}

	return true
}

func calDateTime(name string) (time.Time, bool) {
	nLen := len(name)
	if nLen != 10 && nLen != 19 {
		return time.Now(), false
	}

	tNow := time.Now()
	var yN, mN, dN int
	// 固定的时间，避免每次都对同一张相片产生不同的时间戳
	hN := 15
	minN := 11
	sN := 16

	if nLen == 19 { //2021:01:30 23:36:27
		yStr := name[:4]
		mStr := name[5:7]
		dStr := name[8:10]
		hStr := name[11:13]
		minStr := name[14:16]
		sStr := name[17:19]

		yN, _ = strconv.Atoi(yStr)
		mN, _ = strconv.Atoi(mStr)
		dN, _ = strconv.Atoi(dStr)
		hN, _ = strconv.Atoi(hStr)
		minN, _ = strconv.Atoi(minStr)
		sN, _ = strconv.Atoi(sStr)
	} else if nLen == 10 { // 2021:01:30
		yStr := name[:4]
		mStr := name[5:7]
		dStr := name[8:10]

		yN, _ = strconv.Atoi(yStr)
		mN, _ = strconv.Atoi(mStr)
		dN, _ = strconv.Atoi(dStr)
	} else {
		return tNow, false
	}

	if isValidData(yN, mN, dN, hN, minN, sN) {
		return time.Date(yN, time.Month(mN), dN, hN, minN, sN, 0, time.Local), true
	} else {
		return tNow, false
	}
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

	cntRname++
	if cntRname%100 == 0 {
		log.Printf("cntRname: %d\n", cntRname)
	}
}

func hdlRename2Time(dir, oName string, nTime time.Time) {
	cntMtx.Lock()
	defer cntMtx.Unlock()

	sufIdx := strings.LastIndex(oName, ".") //a.b
	if sufIdx == -1 {
		return
	}
	nName := fmt.Sprintf("%4d%02d%02d%02d%02d%02d%s",
		nTime.Year(), nTime.Month(), nTime.Day(),
		nTime.Hour(), nTime.Minute(), nTime.Second(), oName[sufIdx:])
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
