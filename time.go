package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
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

// IMG_20230204_180346.jpg
// VID_20230101_173306.mp4
func huaweiPhoneTime(path, name string) (retT time.Time, retF bool) {
	if strings.HasPrefix(name, "IMG_") && strings.HasSuffix(name, ".jpg") {
	} else if strings.HasPrefix(name, "VID_") && strings.HasSuffix(name, ".mp4") {
	} else {
		return
	}

	if len(name) < len("VID_20230101_173306.mp4") {
		return
	}

	yStr := name[4:8]
	mStr := name[8:10]
	dStr := name[10:12]
	hStr := name[13:15]
	minStr := name[15:17]
	sStr := name[17:19]

	yN, _ := strconv.Atoi(yStr)
	mN, _ := strconv.Atoi(mStr)
	dN, _ := strconv.Atoi(dStr)
	hN, _ := strconv.Atoi(hStr)
	minN, _ := strconv.Atoi(minStr)
	sN, _ := strconv.Atoi(sStr)

	if isValidData(yN, mN, dN, hN, minN, sN) {
		retT = time.Date(yN, time.Month(mN), dN, hN, minN, sN, 0, time.Local)
		retF = true
	}
	return
}

func shotTimeJPG(path, str string) (retT time.Time, retF bool) {
	allName := fmt.Sprintf("%s\\%s", path, str)

	defer func() {
		if p := recover(); p != nil {
			err := fmt.Errorf("parser %s error: %v", allName, p)
			log.Println(err.Error())
		}
	}()

	exifMtx.Lock()
	defer exifMtx.Unlock()

	if strings.HasSuffix(str, ".jpg") == false {
		if strings.HasSuffix(str, ".JPG") == false {
			return
		}
	}

	f, err := os.Open(allName)
	defer f.Close()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	//exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		//log.Printf(err.Error())
		return
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
		return
	}

	retT = tm
	retF = true
	return
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
		if isValidData(yN, mN, dN, hN, minN, sN) == false {
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

func earlyTime(path, name string) (retT time.Time, retF bool) {
	var t0 time.Time

	t1, ok1 := fileSysTime(path, name)
	t2, ok2 := fileNameTime(path, name)

	if ok1 && ok2 {
		if t1.Before(t2) {
			t0 = t1
		} else {
			t0 = t2
		}
		retT = t0
		retF = true
		return
	} else if ok1 {
		t0 = t1
		retT = t0
		retF = true
		return
	}

	return
}
