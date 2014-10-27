package server

import (
	"github.com/LapisBlue/Tar/head"
	"github.com/LapisBlue/Tar/skin"
	"image/png"
	"log"
	"net/http"
	"strings"
)

func Start(address string) error {
	return http.ListenAndServe(address, &headHandler{})
}

type headHandler struct{}

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

	img, err := head.Render(sk, 45, 256, 256, 4, true, true, true)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "image/png")

	err = png.Encode(w, img)
	if err != nil {
		log.Println(err)
	}
}
