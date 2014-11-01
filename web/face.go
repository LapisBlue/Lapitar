package web

import (
	"github.com/LapisBlue/Lapitar/face"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"net/url"
)

func serveFace(c web.C, w http.ResponseWriter, r *http.Request, params url.Values) {
	watch := util.StartedWatch()

	conf := new(faceConfig)
	*conf = *defaults.Face
	decoder.Decode(conf, params) // Load the settings from the query

	player := c.URLParams["player"]
	sk, err := downloadSkin(player, watch)
	if err != nil {
		return
	}

	watch.Mark()
	result := face.Render(sk, conf.Size, conf.Helm)
	log.Println("Rendered face:", player, watch)

	sendResult(w, player, result, watch)
	watch.Stop()
}

func serveFaceNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, r.URL.Query())
}

func serveFaceWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("size", c.URLParams["size"])
	serveFace(c, w, r, query)
}
