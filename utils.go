package main

import (
	"fmt"
	"log"
	"os"
)

func hdlRename(dir, oName, nName string) {
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
}

func rmFile(path1, path2 string) {
	if DEBUG_MODEL {
		log.Printf("remove  %s  samewith----> %s\n", path1, path2)
	} else {
		os.Remove(path1)
	}
}
