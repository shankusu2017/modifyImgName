package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	formats[0].Head = "微信图片_"
	formats[0].Name = "20240226153256"
	formats[1].Head = "WeChat_"
	formats[1].Name = "20240226153302"
	formats[2].Head = "export"
	formats[2].Name = "1709625041"
}

func main() {
	pwd := "D:\\img\\unsync\\24"
	pwd, _ = os.Getwd()
	pwd = fmt.Sprintf("%s\\test", pwd)
	hdlSame(pwd)
	hdlName(pwd)
	log.Printf("cntRname: %d, cntRemove: %d\n", cntRname, cntRemove)
}
