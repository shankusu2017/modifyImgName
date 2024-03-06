package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	cntRname  = 0
	cntRemove = 0
	cntMtx    sync.RWMutex
)

func hdlRename(dir, oName, nName string) {
	cntMtx.Lock()
	defer cntMtx.Unlock()

	if oName == nName {
		return
	}
	oPath := fmt.Sprintf("%s\\%s", dir, oName)
	nPath := fmt.Sprintf("%s\\%s", dir, nName)
	if DEBUG_MODEL {
		log.Printf("rename  %s ----> %s\n", oPath, nPath)
	} else {
		os.Rename(oPath, nPath)
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
