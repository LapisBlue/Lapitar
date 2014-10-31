package web

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"net/url"
)

func face(w http.ResponseWriter, r *http.Request, params url.Values) {
	conf := new(faceConfig)
	*conf = *defaults.Face
	decoder.Decode(conf, params) // Load the settings from the query
	// TODO: Render face
}

func serveFace(c web.C, w http.ResponseWriter, r *http.Request) {
	face(w, r, r.URL.Query())
}

func serveFaceWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("size", c.URLParams["size"])
	face(w, r, query)
}
