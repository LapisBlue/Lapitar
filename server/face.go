package server

import (
	"github.com/LapisBlue/Lapitar/face"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
)

func serveFace(c web.C, w http.ResponseWriter, r *http.Request, size int, overlay bool) {
	watch := util.StartedWatch()

	conf := defaults.Face
	if size < face.MinimumSize {
		size = face.MinimumSize
	} else if size > conf.Size.Max {
		size = conf.Size.Max
	}

	player := c.URLParams["player"]
	/*meta := loadSkinMeta(player, watch)

	// Check if we can return 304 Not Modified
	if serveCached(w, r, meta) {
		return
	}*/

	skin := downloadSkin(player, watch)
	prepareResponse(w, r)

	watch.Mark()
	result := face.Render(skin, size, overlay, conf.Scale.Get())
	log.Println("Rendered face:", player, watch)

	sendResult(w, player, result, watch)
	watch.Stop()
}

func serveFaceNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, defaults.Head.Size.Def, false)
}

func serveFaceWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, parseSize(c, defaults.Face.Size.Def), false)
}

func serveHelmNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, defaults.Head.Size.Def, true)
}

func serveHelmWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, parseSize(c, defaults.Face.Size.Def), true)
}
