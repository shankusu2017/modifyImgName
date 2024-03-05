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
	mtx sync.Mutex
)

func shotTimeJPG(path, str string) (string, bool) {
	allName := fmt.Sprintf("%s//%s", path, str)
	mtx.Lock()
	defer mtx.Unlock()

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
