package skin

import (
	"errors"
	"fmt"
	"image/png"
	"net/http"
)

const (
	skinURL = "http://skins.minecraft.net/MinecraftSkins/%s.png"
)

func Download(player string) (skin *Skin, err error) {
	resp, err := http.Get(fmt.Sprintf(skinURL, player))
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Request.URL.String() + " returned " + resp.Status)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "image/png" {
		err = errors.New("expected image/png, " + resp.Request.URL.String() + " returned " + contentType + " instead")
		return
	}

	defer resp.Body.Close()
	img, err := png.Decode(resp.Body)
	if err != nil {
		return
	}

	return (*Skin)(rgba(img)), nil
}
