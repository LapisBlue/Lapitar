package cache

import (
	"github.com/LapisBlue/Lapitar/mc"
	"time"
)

var skinCache SkinCache

type SkinCache interface {
	Fetch(uuid string) (SkinMeta, error)
	FetchByName(name string) (SkinMeta, error)
}

type SkinMeta interface {
	mc.SkinMeta
	Profile() mc.Profile
	Timestamp() time.Time
}

func Init(cache SkinCache) {
	skinCache = cache
}

func Fetch(uuid string) (SkinMeta, error) {
	return skinCache.Fetch(uuid)
}

func FetchByName(name string) (SkinMeta, error) {
	return skinCache.FetchByName(name)
}
