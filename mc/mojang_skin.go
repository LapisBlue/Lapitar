package mc

import (
	"fmt"
	"github.com/LapisBlue/Lapitar/util/lhttp"
	"image/png"
)

const (
	LegacyURL = "http://skins.minecraft.net/MinecraftSkins/%s.png"
	// TODO: skinServer = "texture.minecraft.net"
	Steve = "https://minecraft.net/images/steve.png"
	Alex  = "https://minecraft.net/images/alex.png"
)

func SkinURL(player string) string {
	return fmt.Sprintf(LegacyURL, player)
}

func DownloadSkin(url string, alex bool) (sk Skin, err error) {
	req, err := lhttp.Get(url)
	if err != nil {
		return
	}

	req.Header.Set("Accept", lhttp.TypePNG)

	resp, err := lhttp.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = lhttp.ExpectSuccess(resp)
	if err != nil {
		return
	}

	err = lhttp.ExpectContent(resp, lhttp.TypePNG)
	if err != nil {
		return
	}

	img, err := png.Decode(resp.Body)
	if err != nil {
		return
	}

	sk = CreateSkin(img, alex)
	return
}

func (meta mojangSkinMeta) Download() (sk Skin, err error) {
	return DownloadSkin(meta.url, meta.alex)
}
