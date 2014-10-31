package web

import (
	"github.com/LapisBlue/Tar/skin"
	"github.com/LapisBlue/Tar/util"
	"image"
	"image/png"
	"log"
	"net/http"
)

func printError(err error, message ...interface{}) {
	if err == nil {
		return
	}

	log.Println(message...)
	log.Println(err)
}

func downloadSkin(player string, watch *util.StopWatch) (sk *skin.Skin, err error) {
	watch.Mark()
	sk, err = skin.Download(player)
	if err == nil {
		log.Println("Downloaded skin:", player, watch)
	} else {
		printError(err, "Failed to download skin:", player, watch)
	}

	return
}

func sendResult(w http.ResponseWriter, player string, result image.Image, watch *util.StopWatch) (err error) {
	watch.Mark()
	w.Header().Add("Content-Type", "image/png")
	err = png.Encode(w, result)
	if err == nil {
		log.Println("Sent response:", player, watch)
	} else {
		printError(err, "Failed to send response:", player, watch)
	}

	return
}
