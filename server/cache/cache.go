package cache

import "os"

var (
	cacheFolder string
)

func Init(dir string) {
	cacheFolder = dir
	os.MkdirAll(dir, os.ModePerm)
	initSkins()
}
