package server

import (
	"github.com/LapisBlue/Lapitar/head"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
)

func serveHead(c web.C, w http.ResponseWriter, r *http.Request, size int) {
	watch := util.StartedWatch()

	conf := defaults.Head
	if size < head.MinimalSize {
		size = head.MinimalSize
	} else if size > conf.Size.Max {
		size = conf.Size.Max
	}

	player := c.URLParams["player"]
	sk, err := downloadSkin(player, watch)
	if err != nil {
		return
	}

	watch.Mark()
	result, err := head.Render(sk, conf.Angle, size, size, conf.SuperSampling, conf.Helm, conf.Shadow, conf.Lighting, conf.Scale.Get())
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
