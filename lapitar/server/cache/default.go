package cache

import (
	"github.com/LapisBlue/Lapitar/lapitar/mc"
	"time"
)

var (
	steve = downloadDefaultSkin(mc.Steve())
	alex  = downloadDefaultSkin(mc.Alex())
)

func Steve() SkinMeta {
	return steve
}

func Alex() SkinMeta {
	return alex
}

func downloadDefaultSkin(meta mc.SkinMeta, err error) *defaultSkinMeta {
	if err != nil {
		panic("Failed to load default skin: " + err.Error())
	}
	return &defaultSkinMeta{meta, time.Now()}
}

type defaultSkinMeta struct {
	mc.SkinMeta
	timestamp time.Time
}

func (meta *defaultSkinMeta) Profile() mc.Profile {
	panic("Unsupported for default skins")
}

func (meta *defaultSkinMeta) Timestamp() time.Time {
	return meta.timestamp
}

func (meta *defaultSkinMeta) Fetch() (SkinMeta, mc.Skin) {
	sk, err := meta.Download()
	if err != nil {
		panic(err) // This should never return an error
	}
	return meta, sk
}
