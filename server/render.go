package server

import (
	"github.com/LapisBlue/Lapitar/render"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"math"
)

func serveRender(c web.C, w http.ResponseWriter, r *http.Request, size int, conf *renderConfig, portrait, full bool) {
	watch := util.StartedWatch()

	if size < render.MinimalSize {
		size = render.MinimalSize
	} else if size > conf.Size.Max {
		size = conf.Size.Max
	}

	player := c.URLParams["player"]
	meta := loadSkinMeta(player, watch)

	// Check if we can return 304 Not Modified
	if serveCached(w, r, meta) {
		return
	}

	skin := downloadSkin(meta, watch)
	prepareResponse(w, r, skin)

	watch.Mark()
	sizeY := size
	if full {
		sizeY = int(math.Floor(float64(sizeY)*1.625))
	}
	result, err := render.Render(skin.Skin(), conf.Angle, conf.Tilt, conf.Zoom, size, sizeY, conf.SuperSampling, portrait, full, conf.Overlay, conf.Shadow, conf.Lighting, conf.Scale.Get())
	if err == nil {
		log.Println("Rendered head:", player, watch)
	} else {
		printError(err, "Failed to render head:", player, watch)
		return
	}

	sendResult(w, player, result, watch)
	watch.Stop()
}

func serveHeadNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, defaults.Head.Size.Def, defaults.Head, false, false)
}

func serveHeadWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, parseSize(c, defaults.Head.Size.Def), defaults.Head, false, false)
}

func servePortraitNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, defaults.Portrait.Size.Def, defaults.Portrait, true, false)
}

func servePortraitWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, parseSize(c, defaults.Portrait.Size.Def), defaults.Portrait, true, false)
}

func servePlayerNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, defaults.Body.Size.Def, defaults.Body, false, true)
}

func servePlayerWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, parseSize(c, defaults.Body.Size.Def), defaults.Body, false, true)
}
