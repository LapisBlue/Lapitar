package web

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"net/url"
)

func head(w http.ResponseWriter, r *http.Request, params url.Values) {
	conf := new(headConfig)
	*conf = *defaults.Head
	decoder.Decode(conf, params) // Load the settings from the query
	// TODO: Render head
}

func serveHead(c web.C, w http.ResponseWriter, r *http.Request) {
	head(w, r, r.URL.Query())
}

func serveHeadWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("size", c.URLParams["size"])
	head(w, r, query)
}
