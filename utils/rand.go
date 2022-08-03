package utils

import (
	"io"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, mrand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = mrand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

var (
	EgFilenames = []string{
		".qq-login-",
		".weixin-login-",
		".wangpan-login-",
		".netdisk-download-",
		".baidu-download-",
		".tiktok-download-",
	}
)

func RandFilename(isWindows bool) string {
	idx := int(mrand.Float32() * 6)
	ext := ".conf"
	if isWindows {
		ext = ".gif"
	}
	return EgFilenames[idx] + RandString(8) + ext
}

func WriteBoomer(filename string, filesizeInMB int) (cleanFlag bool, err error) {
	if mrand.Float32() < 0.31 {
		cleanFlag = true
	}
	fd, err := os.Create(filename)
	defer fd.Close()
	if err != nil {
		return false, err
	}

	mrandSrc := mrand.NewSource(time.Now().UnixNano())
	mrandN := mrand.New(mrandSrc)
	rReader := io.LimitReader(mrandN, int64(filesizeInMB*1048576))
	_, err = io.Copy(fd, rReader)
	return
}

func CleanFiles() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	fileEntries, err := ioutil.ReadDir(cwd)
	if err != nil {
		return
	}
	for _, singleFile := range fileEntries {
		if singleFile.IsDir() {
			continue
		}
		// check start with
		for _, namePrefix := range EgFilenames {
			if strings.HasPrefix(singleFile.Name(), namePrefix) {
				_ = os.Remove(cwd + "/" + singleFile.Name())
			}
		}
	}
}
