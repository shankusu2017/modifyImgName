package main

import (
	"log"
	"os"
)

func main() {
	pwd, _ := os.Getwd()
	//pwd = fmt.Sprintf("%s\\test", pwd)
	//pwd := "D:\\img\\unsync\\24"
	//pwd = "X:\\img\\"
	pwd = "X:\\img"
	//pwd = "E:\\git\\modifyImgName\\test"
	//hdlSame(pwd)
	hdlName(pwd)
	log.Printf("cntRname: %d, cntRemove: %d\n", cntRname, cntRemove)
}
