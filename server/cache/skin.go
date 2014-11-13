package cache

import (
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/LapisBlue/Lapitar/server/httputil"
	"github.com/LapisBlue/Lapitar/util"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	updateInterval = 24 * time.Hour
)

var (
	skinFolder string
	SteveSkin  Skin
	AlexSkin   Skin

	textureServer = "textures.minecraft.net"
	textureURL    = "http://" + textureServer + "/texture/"
)

func initSkins() {
	skinFolder = filepath.Join(cacheFolder, "skins")
	os.MkdirAll(skinFolder, os.ModePerm)

	SteveSkin = GetSkin("steve").Load()
	AlexSkin = GetSkin("alex").Load()
}

// Basic information about a cached skin
type Meta interface {
	Name() string
	ID() string
	LastMod() time.Time
	Load() Skin
}

// A downloaded skin
type Skin interface {
	Meta
	Skin() *mc.Skin
}

type skinSource struct {
	skinMeta
	path *string
	req  *http.Request
}

func (skin *skinSource) Load() (resultSkin Skin) {
	def := false
	switch skin.id {
	case "steve":
		if SteveSkin != nil {
			return SteveSkin
		}

		def = true // We still need to load Steve
	case "alex":
		if AlexSkin != nil {
			return AlexSkin
		}

		def = true // We still need to load Alex
	}

	result := &localSkin{skinMeta: skin.skinMeta}
	resultSkin = result

	var err error
	defer func() { // Error handler
		if err != nil {
			if def {
				panic(err)
			}

			defSkin := DefaultSkin()
			result.id = defSkin.ID()
			result.skin = defSkin.Skin()
		}
	}()

	var in io.ReadCloser

	if skin.path != nil {
		// We have this skin loaded on the disk
		in, err = os.Open(*skin.path)
		if err != nil {
			in = nil // I'm not sure why I need this but it is not working if I don't do this

			if !os.IsNotExist(err) {
				log.Println(err)
				return
			}

			// We know which skin ID to download, but we haven't downloaded it yet
			if skin.req == nil { // This is not nil if we're downloading the default skins
				skin.req, err = httputil.Get(textureURL + skin.id)
				if err != nil {
					return
				}
				prepareSkinRequest(skin.req)
			}
		} else {
			defer in.Close()
		}
	}

	if skin.req != nil && in == nil {
		// Send the request and download the skin from the server
		resp, err := httputil.Do(skin.req)
		if err != nil {
			return
		}
		in = resp.Body
		defer in.Close()

		// Check if status code was successful
		if !httputil.IsSuccess(resp) {
			err = httputil.NewError(resp, "Expected OK, got "+resp.Status+" instead")
			return
		}

		// Check if response is really a PNG image
		if respType := resp.Header.Get("Content-Type"); respType != "image/png" {
			err = httputil.NewError(resp, "Expected content type image/png, got "+respType+" instead")
			return
		}

		defer func() {
			// If nothing bad happened we still need to save the skin to the cache
			if err == nil {
				file, err := os.Create(filepath.Join(skinFolder, skin.id))
				if err != nil {
					log.Println(err)
					return
				}
				defer file.Close()

				err = png.Encode(file, result.skin.Image())
				if err != nil {
					log.Println(err)
				}
			}
		}()
	}

	if in == nil {
		panic("Invalid skin meta")
	}

	// Decode the PNG image
	img, err := png.Decode(in)
	if err != nil {
		return
	}

	result.skin = mc.CreateSkin(img)
	return
}

func DefaultSkin() Skin {
	if rand.Intn(2) == 0 {
		return SteveSkin
	} else {
		return AlexSkin
	}
}

func GetSkin(name string) Meta {
	if !mc.IsName(name) {
		return DefaultSkin()
	}

	name = mc.ToLower(name)
	meta := new(skinSource)
	meta.name = name

	switch name {
	case "steve", "char":
		// Maybe we need to load steve first
		if SteveSkin == nil {
			return getDefault(meta, "steve", mc.Steve)
		}

		return SteveSkin
	case "alex":
		// Maybe we need to load steve first
		if AlexSkin == nil {
			return getDefault(meta, "alex", mc.Alex)
		}

		return AlexSkin
	}

	// Check if the skin is cached on disk
	err := loadCachedMeta(meta)
	if err == nil {
		// Check if we need to update the skin
		if time.Now().Sub(meta.LastMod()) > updateInterval {
			update := new(skinSource)
			*update = *meta
			go func() {
				watch := util.StartedWatch()
				querySkinMeta(update)
				if meta.id != update.id {
					// The skin is outdated, we need to update it
					updateSkinMeta(meta)
					// Delete the old skin
					err := os.Remove(*meta.path)
					if err != nil {
						log.Println(err)
					}
				} else {
					updateSkinMeta(meta)
				}

				log.Println("Checked for skin update:", meta.name, watch)
			}()
		}

		return meta
	} else if !os.IsNotExist(err) {
		log.Println(err)
	}

	// We don't have this skin in memory, we need to query the Mojang server about the skin
	err = querySkinMeta(meta)
	if err == nil {
		linkSkinMeta(meta)
	} else {
		log.Println(err)

		// Assign the default skin to the name
		def := DefaultSkin()

		skin := &localSkin{skinMeta: meta.skinMeta}

		skin.id = def.Name()
		skin.lastMod = time.Now()
		skin.skin = def.Skin()
		linkSkinMeta(skin)
		return skin
	}

	return meta
}

func getDefault(meta *skinSource, name string, source string) Meta {
	meta.name = name
	meta.id = name

	err := loadDefaultMeta(meta)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(name + " not found, downloading...")
		} else {
			log.Println(err)
		}
	}

	meta.req, err = httputil.Get(source)
	if err != nil {
		panic(err)
	}

	prepareSkinRequest(meta.req)
	return meta
}

func loadDefaultMeta(meta *skinSource) (err error) {
	path := filepath.Join(skinFolder, meta.name)
	stat, err := os.Stat(path)
	if err != nil {
		return
	}

	meta.lastMod = stat.ModTime()
	meta.path = &path
	return
}

func loadCachedMeta(meta *skinSource) (err error) {
	path := filepath.Join(skinFolder, meta.name)
	stat, err := os.Lstat(path)
	if err != nil {
		return
	}

	target, err := os.Readlink(path)
	if err != nil {
		return
	}

	meta.lastMod = stat.ModTime()
	meta.id = filepath.Base(target)
	meta.path = &target
	return
}

func prepareSkinRequest(req *http.Request) {
	// We only accept PNG images as response
	req.Header.Set("Accept", "image/png")
}

func querySkinMeta(meta *skinSource) (err error) {
	req, err := httputil.Get(mc.SkinURL(meta.name))
	if err != nil {
		return
	}

	prepareSkinRequest(req)

	// Try to get the skin ID from textures.minecraft.net
	loc, err := httputil.GetLocation(req, textureServer)
	if err != nil {
		return
	}

	meta.id = filepath.Base(loc.URL.Path)
	meta.req = loc
	return
}

func linkSkinMeta(meta Meta) {
	if err := os.Symlink(meta.ID(), filepath.Join(skinFolder, meta.Name())); err != nil {
		log.Println(err)
	}
}

func updateSkinMeta(meta Meta) {
	link := filepath.Join(skinFolder, meta.Name())
	err := os.Remove(link)
	if err != nil {
		log.Println(err)
	}

	linkSkinMeta(meta)
}

// struct implementation

type skinMeta struct {
	name    string
	id      string
	lastMod time.Time
}

func (meta skinMeta) Name() string {
	return meta.name
}

func (meta skinMeta) ID() string {
	return meta.id
}

func (meta skinMeta) LastMod() time.Time {
	return meta.lastMod
}

type localSkin struct {
	skinMeta
	skin *mc.Skin
}

func (skin *localSkin) Load() Skin {
	return skin
}

func (skin *localSkin) Skin() *mc.Skin {
	return skin.skin
}
