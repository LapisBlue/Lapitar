package server

import (
	"github.com/LapisBlue/Lapitar/face"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
)

func serveFace(c web.C, w http.ResponseWriter, r *http.Request, size int) {
	watch := util.StartedWatch()

	conf := defaults.Face
	if size < face.MinimalSize {
		size = face.MinimalSize
	} else if size > conf.Size.Max {
		size = conf.Size.Max
	}

	player := c.URLParams["player"]
	sk, err := downloadSkin(player, watch)
	if err != nil {
		return
	}

	watch.Mark()
	result := face.Render(sk, size, conf.Helm, conf.Scale.Get())
	log.Println("Rendered face:", player, watch)

	sendResult(w, player, result, watch)
	watch.Stop()
}

func serveFaceNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, defaults.Head.Size.Def)
}

func serveFaceWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, parseSize(c, defaults.Face.Size.Def))
}
