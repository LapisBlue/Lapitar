package server

import (
	"github.com/LapisBlue/Lapitar/render"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
)

func serveHead(c web.C, w http.ResponseWriter, r *http.Request, size int) {
	watch := util.StartedWatch()

	conf := defaults.Head
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
	result, err := render.RenderHead(skin.Skin(), conf.Angle, 20, -4.5, size, size, conf.SuperSampling, conf.Helm, conf.Shadow, conf.Lighting, conf.Scale.Get())
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
	serveHead(c, w, r, defaults.Head.Size.Def)
}

func serveHeadWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveHead(c, w, r, parseSize(c, defaults.Head.Size.Def))
}
