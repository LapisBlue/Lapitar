package skin

import (
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"regexp"
)

const (
	skinURL = "http://skins.minecraft.net/MinecraftSkins/%s.png"
)

var (
	namePattern = regexp.MustCompile("^[a-zA-Z0-9_]{1,16}$")
	steve       = "https://minecraft.net/images/steve.png"
	alex        = "https://minecraft.net/images/alex.png" // TODO
)

func Download(player string) (skin *Skin, err error) {
	var url string
	if namePattern.MatchString(player) {
		url = fmt.Sprintf(skinURL, player)
	} else {
		url = steve
	}

	resp, err := http.Get(url)
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

	return Create(img), nil
}
