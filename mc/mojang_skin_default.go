package mc

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

const (
	steveURL = "https://minecraft.net/images/steve.png"
	alexURL  = "https://minecraft.net/images/alex.png"
)

func downloadDefaultSkin(meta *defaultSkinMeta, alex bool) (err error) {
	resp, err := requestSkin(meta.url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	hash := md5.New()
	reader := io.TeeReader(resp.Body, hash)
	meta.skin, err = readSkin(reader, alex)
	if err != nil {
		return
	}

	meta.hash = hex.EncodeToString(hash.Sum(nil))
	return
}

func Steve() (meta SkinMeta, err error) {
	sk := &defaultSkinMeta{url: steveURL}
	err = downloadDefaultSkin(sk, false)
	if err != nil {
		return
	}

	meta = sk
	return
}

func Alex() (meta SkinMeta, err error) {
	sk := &defaultSkinMeta{url: alexURL}
	err = downloadDefaultSkin(sk, true)
	if err != nil {
		return
	}

	meta = sk
	return
}

type defaultSkinMeta struct {
	hash string
	url  string
	skin Skin
}

func (meta *defaultSkinMeta) ID() string {
	return meta.hash
}

func (meta *defaultSkinMeta) URL() string {
	return meta.url
}

func (meta *defaultSkinMeta) Download() (Skin, error) {
	return meta.skin, nil
}
