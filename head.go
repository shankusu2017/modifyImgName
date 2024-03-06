package main

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type nameFormat struct {
	Head string
	Name string
}

var (
	formats [3]nameFormat
)

func delHead(str string) (string, bool) {
	for _, format := range formats {
		ok := strings.HasPrefix(str, format.Head)
		if ok == false {
			continue
		}
		if len(str) < len(format.Head)+len(format.Name) {
			continue
		}

		return str[len(format.Head):], true
	}

	return "", false
}

func del32CharHead(path, str string) (string, bool) {
	//a0aab87020cd2ca6dbe5323bcbc630ee.mp4
	idx := strings.Index(str, ".")
	if idx != 32 {
		return "", false
	}
	_, err := hex.DecodeString(str[:idx])
	if err != nil {
		return "", false
	}

	ftime, ok := fileSysTime(path, str)
	if ok == true {
		nName := fmt.Sprintf("%4d%02d%02d%02d%02d%02d",
			ftime.Year(), ftime.Month(), ftime.Day(),
			ftime.Hour(), ftime.Minute(), ftime.Second())
		return nName + str[32:], true
	} else {
		return "", false
	}
}
