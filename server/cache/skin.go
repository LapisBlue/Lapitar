package cache

import (
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/LapisBlue/Lapitar/server/httputil"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

var (
	skinFolder string
	skinClient = httputil.Client()
	steveSkin  *mc.Skin
	alexSkin   *mc.Skin
)

func init() {
	/*skinClient.RedirectHandler = func(req *http.Request, via []*http.Request) bool {
		if req.Host == "textures.minecraft.net" {
			hash := filepath.Base(req.URL.Path)
			if _, err := os.Stat(filepath.Join(skinFolder, hash)); err == nil {
				return false // We have this skin already
			}
		}
		return false
	}// TODO: do I need this*/
}

func initSkins() {
	skinFolder = filepath.Join(cacheFolder, "skins")
	os.MkdirAll(skinFolder, perms)

	var err error
	steveSkin, err = GetSkin("steve")
	if err != nil {
		panic("Failed to load Steve: " + err.Error())
	}
	alexSkin, err = GetSkin("alex")
	if err != nil {
		panic("Failed to load Alex: " + err.Error())
	}
}

func loadSkinCached(name string) (skin *mc.Skin, err error) {
	if name == "steve" && steveSkin != nil {
		return steveSkin, nil
	} else if name == "alex" && alexSkin != nil {
		return alexSkin, nil
	}

	file, err := os.Open(filepath.Join(skinFolder, name))
	if err != nil {
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return
	}

	skin = mc.CreateSkin(img)
	return
}

func writeSkinCached(name string, player string, skin *mc.Skin) (err error) {
	file, err := os.Create(filepath.Join(skinFolder, name))
	if err != nil {
		return
	}
	defer file.Close()

	err = png.Encode(file, skin.Image())
	if err != nil {
		return
	}

	if name != player {
		err = os.Symlink(name, filepath.Join(skinFolder, player))
	}

	return
}

// TODO: Configurable
func DefaultSkin() (*mc.Skin, error) {
	if rand.Intn(1) == 0 {
		return GetSkin("steve")
	} else {
		return GetSkin("alex")
	}
}

func GetSkin(player string) (skin *mc.Skin, err error) {
	if !mc.IsName(player) {
		return DefaultSkin()
	}

	player = mc.ToLower(player)
	if player == "char" { // Alias for Steve
		player = "steve"
	}

	skin, err = loadSkinCached(player)
	if err == nil {
		return
	}

	def := false

	var url string
	switch player {
	case "steve":
		url = mc.Steve
		def = true
	case "alex":
		url = mc.Alex
		def = true
	default:
		url = mc.SkinURL(player)
	}

	req, err := skinClient.Get(url)
	if err != nil {
		log.Println(err)
		if def {
			return
		} else {
			return DefaultSkin()
		}
	}

	req.Header.Set("Accept", "image/png")

	resp, err := skinClient.Do(req)
	if err != nil {
		log.Println(err)
		if def {
			return
		} else {
			return DefaultSkin()
		}
	}
	defer resp.Body.Close()

	img, err := png.Decode(resp.Body)
	if err != nil {
		log.Println(err)
		if def {
			return
		} else {
			return DefaultSkin()
		}
	}

	skin = mc.CreateSkin(img)

	name := player
	if !def {
		name = filepath.Base(resp.Request.URL.Path) // The texture ID
	}

	err = writeSkinCached(name, player, skin)
	if err != nil {
		log.Println(err)
	}

	return
}
