package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	pwd := "D:\\img\\unsync\\24"
	pwd, _ = os.Getwd()
	pwd = fmt.Sprintf("%s\\test", pwd)
	hdlSame(pwd)
	hdlName(pwd)
	log.Printf("cntRname: %d, cntRemove: %d\n", cntRname, cntRemove)
}
