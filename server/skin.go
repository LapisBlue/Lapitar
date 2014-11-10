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

	skin := downloadSkin(meta, watch)
	prepareResponse(w, r, skin)

	sendResult(w, player, skin.Skin().Image(), watch)
	watch.Stop()
}
