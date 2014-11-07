package cache

import (
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/LapisBlue/Lapitar/server/httputil"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	steveSkin, _, err = GetSkin("steve")
	if err != nil {
		panic("Failed to load Steve: " + err.Error())
	}
	alexSkin, _, err = GetSkin("alex")
	if err != nil {
		panic("Failed to load Alex: " + err.Error())
	}
}

func loadSkinCached(name string) (skin *mc.Skin, id string, duration time.Duration, err error) {
	id = name

	if name == "steve" && steveSkin != nil {
		skin = steveSkin
		return
	} else if name == "alex" && alexSkin != nil {
		skin = alexSkin
		return
	}

	path := filepath.Join(skinFolder, name)

	if target, err := os.Readlink(path); err == nil {
		id = target
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	duration = time.Now().Sub(stat.ModTime())

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
func DefaultSkin() (*mc.Skin, string, error) {
	if rand.Intn(2) == 0 {
		return GetSkin("steve")
	} else {
		return GetSkin("alex")
	}
}

func GetSkin(player string) (skin *mc.Skin, id string, err error) {
	if !mc.IsName(player) {
		return DefaultSkin()
	}

	player = mc.ToLower(player)
	if player == "char" { // Alias for Steve
		player = "steve"
	}

	skin, id, _, err = loadSkinCached(player)
	if err == nil {
		return
	}

	id = player

	var url string
	def := true

	switch player {
	case "steve":
		url = mc.Steve
	case "alex":
		url = mc.Alex
	default:
		url = mc.SkinURL(player)
		def = false
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

	if resp.StatusCode != http.StatusOK {
		err = httputil.NewError(resp, "Expected OK, got "+resp.Status+" instead")
		log.Println(err)
		if def {
			return
		} else {
			return DefaultSkin()
		}
	}

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
		id = name
	}

	err = writeSkinCached(name, player, skin)
	if err != nil {
		log.Println(err)
	}

	return
}
