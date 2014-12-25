package server

import (
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"net/http"
)

func serveSkin(c web.C, w http.ResponseWriter, r *http.Request) {
	watch := util.StartedWatch()

	player := c.URLParams["player"]
	meta := loadSkinMeta(player, watch)

	// Check if we can return 304 Not Modified
	if serveCached(w, r, meta) {
		return
	}

	skin, err := meta.Download()
	if err != nil {
		panic(err)
	}
	prepareResponse(w, r, meta)

	sendResult(w, player, skin.Image(), watch)
	watch.Stop()
}
