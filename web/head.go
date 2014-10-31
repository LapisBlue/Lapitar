package web

import (
	"github.com/LapisBlue/Tar/head"
	"github.com/LapisBlue/Tar/util"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"net/url"
)

func serveHead(c web.C, w http.ResponseWriter, r *http.Request, params url.Values) {
	watch := util.StartedWatch()

	conf := new(headConfig)
	*conf = *defaults.Head
	decoder.Decode(conf, params) // Load the settings from the query

	player := c.URLParams["player"]
	sk, err := downloadSkin(player, watch)
	if err != nil {
		return
	}

	watch.Mark()
	result, err := head.Render(sk, conf.Angle, conf.Size, conf.Size, conf.SuperSampling, conf.Helm, conf.Shadow, conf.Lighting)
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
	serveHead(c, w, r, r.URL.Query())
}

func serveHeadWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("size", c.URLParams["size"])
	serveHead(c, w, r, query)
}
