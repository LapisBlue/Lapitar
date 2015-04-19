package cache

import (
	"github.com/LapisBlue/lapitar/mc"
	"time"
)

var (
	steve = downloadSteve()
	alex  = downloadAlex()
)

func Steve() SkinMeta {
	return steve
}

func Alex() SkinMeta {
	return alex
}

func downloadSteve() *defaultSkinMeta {
	meta, err := mc.Steve()
	return downloadDefaultSkin("Steve", meta, err)
}

func downloadAlex() *defaultSkinMeta {
	meta, err := mc.Alex()
	return downloadDefaultSkin("Alex", meta, err)
}

func downloadDefaultSkin(name string, meta mc.SkinMeta, err error) *defaultSkinMeta {
	if err != nil {
		panic("Failed to load default skin: " + err.Error())
	}
	result := &defaultSkinMeta{SkinMeta: meta, timestamp: time.Now(), name: name}
	_, sk := result.Fetch()
	result.alex = sk.IsAlex()
	return result
}

type defaultSkinMeta struct {
	mc.SkinMeta
	timestamp time.Time
	name      string
	alex      bool
}

func (meta *defaultSkinMeta) Profile() mc.Profile {
	return meta
}

func (meta *defaultSkinMeta) Name() string {
	return meta.name
}

func (meta *defaultSkinMeta) UUID() string {
	panic("Unsupported for default skins")
}

func (meta *defaultSkinMeta) IsAlex() bool {
	return meta.alex
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
