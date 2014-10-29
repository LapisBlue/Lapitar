package web

import (
	"github.com/LapisBlue/Tar/head"
	"github.com/LapisBlue/Tar/skin"
	"github.com/LapisBlue/Tar/util"
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

	log.Println("Processing", player, "requested by", addrFor(r))
	watch := util.CreateStopWatch()

	watch.Start()
	sk, err := skin.Download(player)
	watch.Stop()
	if err != nil {
		log.Println(err, watch)
		return
	}
	log.Println("Downloaded skin:", player, watch)

	watch.Reset().Start()
	img, err := head.Render(sk, 45, 256, 256, 4, true, true, true)
	watch.Stop()
	if err != nil {
		log.Println(err, watch)
		return
	}
	log.Println("Rendered head:", player, watch)

	watch.Reset().Start()
	w.Header().Add("Content-Type", "image/png")

	err = png.Encode(w, img)
	watch.Stop()
	if err != nil {
		log.Println(err, watch)
		return
	}

	log.Println("Response prepared:", player, watch)
}

func addrFor(r *http.Request) (addr string) {
	addr = r.RemoteAddr
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		addr += " (" + forward + ")"
	}
	return
}
