package utils

import (
	"io/ioutil"
	"os"

	"github.com/shirou/gopsutil/v3/disk"
)

func RetrieveSpareSpace() (freeSize uint64, workMode int, writable bool) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// detect current disk free space
	ustat, err := disk.Usage(cwd)
	if err != nil {
		panic(err)
	}
	if ustat.UsedPercent > 99.8 {
		return 0, -1, false
	}
	freeSize = ustat.Free / uint64(1048576) // in MiB
	// detect if current working folder is writable, if true, next
	if freeSize > 2048 {
		workMode = 3
	} else if freeSize >= 64 {
		workMode = 2
	} else {
		workMode = -1
		return
	}
	// try write a file
	tmpFile := cwd + "/." + RandString(6) + ".txt"
	err = ioutil.WriteFile(tmpFile, []byte(RandString(6)), 0644)
	if err != nil {
		return
	} else {
		_ = os.Remove(tmpFile)
		writable = true
		return
	}
}
