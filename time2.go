package main

import (
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"io"
	"os"
	"strings"
	"time"
)

func showTime2JPG(path, name string) (retT time.Time, retF bool) {
	for i := 0; i < EXIF_TRY_TIMES; i++ {
		retT, retF = showTime2JPGHdl(path, name)
		if retF {
			return
		}
	}
	return
}

func showTime2JPGHdl(path, name string) (retT time.Time, retF bool) {
	exifMtx.Lock()
	defer exifMtx.Unlock()

	if strings.HasSuffix(name, ".jpg") == false {
		if strings.HasSuffix(name, ".JPG") == false {
			return
		}
	}

	pathAll := fmt.Sprintf("%s\\%s", path, name)
	retT = time.Now()

	defer func() {
		if errRaw := recover(); errRaw != nil {
			fmt.Sprintf("%v", errRaw)
			return
		}
	}()

	f, err := os.Open(pathAll)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return
	}

	rawExif, err := exif.SearchAndExtractExifN(data, 0)
	if err != nil {
		return
	}
	// Run the parse.
	entries, _, err := exif.GetFlatExifDataUniversalSearch(rawExif, nil, false)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.TagId == 0x9004 ||
			entry.TagId == 0x9003 ||
			entry.TagId == 0x001d ||
			entry.TagId == 0x0132 {
			//fmt.Printf("IFD-PATH=[%s] ID=(0x%04x) NAME=[%s] COUNT=(%d) TYPE=[%s] VALUE=[%s]\n\n", entry.IfdPath, entry.TagId, entry.TagName, entry.UnitCount, entry.TagTypeName, entry.Formatted)
			retTime, ok := calDateTime(entry.Formatted)
			if ok {
				if retTime.Before(retT) {
					retT = retTime
					retF = true
				}
			}
		}
	}
	return
}
