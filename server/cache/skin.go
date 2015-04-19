package cache

import (
	"github.com/LapisBlue/lapitar/mc"
	"time"
)

var skinCache SkinCache

type SkinCache interface {
	Fetch(uuid string) SkinMeta
	FetchByName(name string) SkinMeta
}

type SkinMeta interface {
	mc.SkinMeta
	Profile() mc.Profile
	Timestamp() time.Time
	Fetch() (SkinMeta, mc.Skin)
}

func Init(cache SkinCache) {
	skinCache = cache
}

func Fetch(uuid string) SkinMeta {
	return skinCache.Fetch(uuid)
}

func FetchByName(name string) SkinMeta {
	return skinCache.FetchByName(name)
}
