package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	exifMtx sync.Mutex
)

func shotTimeJPG(path, str string) (string, bool) {
	allName := fmt.Sprintf("%s//%s", path, str)
	exifMtx.Lock()
	defer exifMtx.Unlock()

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

// 从文件系统数据中能读到的最早时间
func fileSysTime(path, str string) (time.Time, bool) {
	tNow := time.Now()
	done := false

	allName := fmt.Sprintf("%s//%s", path, str)
	finfo, err := os.Stat(allName)
	if err != nil {
		return tNow, done
	}
	// 获取文件原来的访问时间，修改时间
	// windows下代码如下
	winFileAttr := finfo.Sys().(*syscall.Win32FileAttributeData)
	ctime := time.Unix(winFileAttr.CreationTime.Nanoseconds()/1e9, 0)
	latime := time.Unix(winFileAttr.LastAccessTime.Nanoseconds()/1e9, 0)
	lwtime := time.Unix(winFileAttr.LastWriteTime.Nanoseconds()/1e9, 0)

	if ctime.IsZero() == false {
		if ctime.Before(tNow) {
			tNow = ctime
			done = true
		}
	}
	if latime.IsZero() == false {
		if latime.Before(tNow) {
			tNow = latime
			done = true
		}
	}
	if lwtime.IsZero() == false {
		if lwtime.Before(tNow) {
			tNow = lwtime
			done = true
		}
	}

	return tNow, done
}

func fileNameTime(path, str string) (time.Time, bool) {
	validEarlyTime := time.Date(2000, 1, 1, 1, 1, 1, 0, time.Local)
	validLasteTime := time.Now()

	// find .
	idx := strings.Index(str, ".")
	if idx == -1 {
		return time.Now(), false
	}
	name := str[:idx]

	// 1709625041.jpg unix时间戳的格式
	// 20240305161544.jpg 文件名格式
	num, err := strconv.Atoi(name)
	if err == nil && int64(num) < validLasteTime.Unix() {
		ut := time.Unix(int64(num), 0)
		if ut.After(validEarlyTime) {
			if ut.Before(validLasteTime) {
				return ut, true
			}
		}
		return time.Now(), false
	} else {
		if len(name) != len("20240305161544") {
			return time.Now(), false
		}
		yStr := name[:4]
		mStr := name[4:6]
		dStr := name[6:8]
		hStr := name[8:10]
		minStr := name[10:12]
		sStr := name[12:14]

		yN, _ := strconv.Atoi(yStr)
		mN, _ := strconv.Atoi(mStr)
		dN, _ := strconv.Atoi(dStr)
		hN, _ := strconv.Atoi(hStr)
		minN, _ := strconv.Atoi(minStr)
		sN, _ := strconv.Atoi(sStr)
		if (mN < 1 || mN > 12) ||
			(dN < 1 || dN > 31) ||
			(hN < 0 || hN > 23) ||
			(minN < 0 || minN > 59) ||
			(sN < 0 || sN > 59) {
			return time.Now(), false
		}
		nTime := time.Date(yN, time.Month(mN), dN, hN, minN, sN, 0, time.Local)
		if nTime.After(validEarlyTime) {
			if nTime.Before(validLasteTime) {
				return nTime, true
			}
		}

		return time.Now(), false
	}

	return time.Now(), false
}

func earlyTime(path, str string) (string, bool) {
	var t0 time.Time
	var n0 string

	idx := strings.LastIndex(str, ".")
	if idx == -1 {
		return "", false
	}

	t1, ok1 := fileSysTime(path, str)
	t2, ok2 := fileNameTime(path, str)

	if ok1 && ok2 {
		if t1.Before(t2) {
			t0 = t1
		} else {
			t0 = t2
		}
		n0 = fmt.Sprintf("%4d%02d%02d%02d%02d%02d%s",
			t0.Year(), t0.Month(), t0.Day(),
			t0.Hour(), t0.Minute(), t0.Second(),
			str[idx:])

		return n0, true
	} else if ok1 {
		t0 = t1
		n0 = fmt.Sprintf("%4d%02d%02d%02d%02d%02d%s",
			t0.Year(), t0.Month(), t0.Day(),
			t0.Hour(), t0.Minute(), t0.Second(),
			str[idx:])

		return n0, true
	}

	return n0, false
}
