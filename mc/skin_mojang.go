package mc

import (
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"strings"
)

const (
	skinURL = "http://skins.minecraft.net/MinecraftSkins/%s.png"
	Steve   = "https://minecraft.net/images/steve.png"
	Alex    = "https://minecraft.net/images/alex.png"
)

func SkinURL(player string) string {
	return fmt.Sprintf(skinURL, player)
}

func Download(player string) (skin *Skin, err error) {
	var url string
	url = fmt.Sprintf(skinURL, player)
	// Char is only supported for compatibility with previous avatar services
	switch {
	case strings.EqualFold(player, "steve"):
		url = Steve
	case strings.EqualFold(player, "alex"):
		url = Alex
	case IsName(player):
		url = fmt.Sprintf(skinURL, player)
	default:
		url = Steve // TODO
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

	return CreateSkin(img), nil
}
