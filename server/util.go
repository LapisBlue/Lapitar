package server

import (
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/zenazn/goji/web"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"
)

func printError(err error, message ...interface{}) {
	if err == nil {
		return
	}

	log.Println(message...)
	log.Println(err)
}

func parseSize(c web.C, def int) (result int) {
	size := c.URLParams["size"]
	result, err := strconv.Atoi(size)
	if err != nil {
		printError(err, "Failed to parse size", size)
		return def
	}
	return
}

func downloadSkin(name string, watch *util.StopWatch) (skin mc.Skin) {
	watch.Mark()
	if !mc.IsUUID(name) {
		if !mc.IsName(name) {
			panic("NO NAME")
		}

		profile, err := mc.FetchProfile(name)
		if err != nil {
			panic(err)
		}
		if profile == nil {
			panic("NO PROFILE")
		}

		name = profile.UUID()
	}

	profile, err := mc.FetchSkin(name)
	if err != nil {
		panic(err)
	}
	if profile == nil {
		panic("NO SKIN")
	}
	skin, err = profile.Skin().Download()
	if err != nil {
		panic(err)
	}

	log.Println("Loaded skin:", name, watch)
	return
}

const (
	keepCache    = 24 * time.Hour
	cacheControl = "max-age=86400" // 24*60*60, one day in seconds
)

func serveCached(w http.ResponseWriter, r *http.Request) bool {
	/*if tag := r.Header.Get("If-None-Match"); tag == meta.ID() {
		prepareResponse(w, r)
		w.WriteHeader(http.StatusNotModified)
		return true
	}*/

	return false
}

func prepareResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", cacheControl)
	w.Header().Add("Expires", time.Now().Add(keepCache).UTC().Format(http.TimeFormat))
	/*w.Header().Add("ETag", meta.ID())
	w.Header().Add("Last-Modified", meta.LastMod().UTC().Format(http.TimeFormat))*/
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
