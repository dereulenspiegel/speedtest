// +build embed

package main

import (
	"log"
	"net/http"

	_ "github.com/dereulenspiegel/speedtest/statik"
	"github.com/rakyll/statik/fs"
)

func serveStatic() http.Handler {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("Failed to load embedded files: %s", err)
	}
	return http.FileServer(statikFS)
}
