package server

import (
	"flag"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"log"
	"net/http"
	"os"
)

var (
	defaults *config
	//decoder  = schema.NewDecoder()
)

func start(conf *config, www string) {
	defaults = conf
	flag.Set("bind", conf.Address) // Uh, I guess that's a bit strange
	if conf.Proxy {
		goji.Insert(middleware.RealIP, middleware.Logger)
	}

	goji.Use(serveLapitar)

	register("/skin/:player", serveSkin)

	register("/face/:player", serveFaceNormal)
	register("/face/:size/:player", serveFaceWithSize)
	register("/helm/:player", serveHelmNormal)
	register("/helm/:size/:player", serveHelmWithSize)

	register("/head/:player", serveHeadNormal)
	register("/head/:size/:player", serveHeadWithSize)
	register("/portrait/:player", servePortraitNormal)
	register("/portrait/:size/:player", servePortraitWithSize)
	register("/body/:player", serveBodyNormal)
	register("/body/:size/:player", serveBodyWithSize)

	if exists(www) {
		goji.Get("/*", http.FileServer(http.Dir(www)))
	} else {
		log.Println("Failed to find website files at", www)
	}

	goji.Serve()
}

func exists(dir string) bool {
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func register(pattern string, handler interface{}) {
	goji.Get(pattern+".png", handler)
	goji.Get(pattern, handler)
}

func serveLapitar(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Lapitar") // TODO: Version
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
