package mc

import (
	"fmt"
	"github.com/LapisBlue/Lapitar/lapitar/util/lhttp"
	"image/png"
	"io"
	"net/http"
)

const (
	LegacyURL = "http://skins.minecraft.net/MinecraftSkins/%s.png"
	// TODO: skinServer = "texture.minecraft.net"
)

func SkinURL(player string) string {
	return fmt.Sprintf(LegacyURL, player)
}

func DownloadSkin(url string, alex bool) (sk Skin, err error) {
	resp, err := requestSkin(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return readSkin(resp.Body, alex)
}

func readSkin(reader io.Reader, alex bool) (sk Skin, err error) {
	img, err := png.Decode(reader)
	if err != nil {
		return
	}

	sk = CreateSkin(img, alex)
	return
}

func requestSkin(url string) (resp *http.Response, err error) {
	req, err := lhttp.Get(url)
	if err != nil {
		return
	}

	req.Header.Set("Accept", lhttp.TypePNG)

	resp, err = lhttp.Do(req)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			resp.Body.Close()
		}
	}()

	err = lhttp.ExpectSuccess(resp)
	if err != nil {
		return
	}

	err = lhttp.ExpectContent(resp, lhttp.TypePNG)
	return
}

func (meta mojangSkinMeta) Download() (sk Skin, err error) {
	return DownloadSkin(meta.url, meta.alex)
}
