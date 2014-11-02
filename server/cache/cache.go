package cache

import "os"

var (
	perms       os.FileMode = os.ModePerm
	cacheFolder string
)

func Init(dir string) {
	cacheFolder = dir
	os.MkdirAll(dir, perms)
	initSkins()
}
