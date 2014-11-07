package server

import (
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"net/http"
)

func serveSkin(c web.C, w http.ResponseWriter, r *http.Request) {
	watch := util.StartedWatch()

	player := c.URLParams["player"]
	sk, id, err := downloadSkin(player, watch)
	if err != nil {
		return
	}

	if !prepareResponse(w, r, id) {
		return
	}

	sendResult(w, player, sk.Image(), watch)
	watch.Stop()
}
