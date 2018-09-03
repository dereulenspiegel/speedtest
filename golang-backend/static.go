// +build !embed

package main

import (
	"flag"
	"net/http"
)

var (
	staticDir = flag.String("static.dir", "../", "Directory to server static files from")
)

func serveStatic() http.Handler {
	return http.FileServer(http.Dir(*staticDir))
}
