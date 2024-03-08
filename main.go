package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	pwd, _ := os.Getwd()
	pwd = fmt.Sprintf("%s\\test", pwd)
	//path := ""
	//name := ""
	pwd = "D:\\img\\unsync\\24"
	//pwd = "X:\\img"
	//pwd = "E:\\git\\modifyImgName\\test"

	//path = "D:\\git\\modifyImgName\\test\\hastime"
	//name = "20210130233627.JPG"
	//name = "20242261340370.jpg"
	//name = "20230909173017.jpg"
	//name = "20230813162626.jpg"

	//showTime2JPG(path, name)
	//log.Printf("=======================\n")
	//nTime, ok := shotTimeJPG(path, name)
	//if ok {
	//	log.Printf(nTime.String())
	//	log.Printf(nTime.String())
	//	log.Printf(nTime.String())
	//}
	if FUNC_DEL_SAME {
		hdlSame(pwd)
	}
	hdlName(pwd)
	log.Printf("cntRname: %d, cntRemove: %d\n", cntRname, cntRemove)
}
