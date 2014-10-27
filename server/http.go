package server

import (
	"github.com/LapisBlue/Tar/head"
	"github.com/LapisBlue/Tar/head/glctx"
	"github.com/LapisBlue/Tar/skin"
	"image/png"
	"log"
	"net/http"
	"strings"
)

func Start(address string) error {
	factory := glctx.OSMesa()
	renderer := &head.Renderer{
		GL: factory,

		Angle:         45,
		Width:         256,
		Height:        256,
		SuperSampling: 4,

		Helmet:   true,
		Shadow:   true,
		Lighting: true,
	}
	return http.ListenAndServe(address, &headHandler{renderer})
}

type headHandler struct {
	r *head.Renderer
}

func (h *headHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.RequestURI
	if player[0] == '/' {
		player = player[1:]
	}
	if strings.HasSuffix(player, ".png") {
		player = player[:len(player)-4]
	}

	log.Println("Downloading skin:", player)
	sk, err := skin.Download(player)
	if err != nil {
		log.Println(err)
	}

	img, err := h.r.Render(sk)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "image/png")

	err = png.Encode(w, img)
	if err != nil {
		log.Println(err)
	}
}
