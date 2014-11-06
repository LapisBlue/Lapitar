package server

import (
	"flag"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"
	"net/http"
)

var (
	defaults *config
	//decoder  = schema.NewDecoder()
)

func start(conf *config) {
	defaults = conf
	flag.Set("bind", conf.Address) // Uh, I guess that's a bit strange
	if conf.Proxy {
		goji.Insert(middleware.RealIP, middleware.Logger)
	}

	register("/skin/:player", serveSkin)

	register("/head/:player", serveHeadNormal)
	register("/head/:size/:player", serveHeadWithSize)

	register("/face/:player", serveFaceNormal)
	register("/face/:size/:player", serveFaceWithSize)

	goji.Get("/*", http.FileServer(http.Dir("www"))) // TODO: How to find the correct dir?

	goji.Serve()
}

func register(pattern string, handler interface{}) {
	goji.Get(pattern+".png", handler)
	goji.Get(pattern, handler)
}
