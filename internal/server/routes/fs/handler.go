// Package fs is the HTTP handler to serve files from an http.FileSystem.
package fs

import (
	"net/http"
	"strings"
)

type handler struct {
	uriPrefix  string
	fileServer http.Handler
}

func NewHandler(uriPrefix string, fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return &handler{
		uriPrefix:  uriPrefix,
		fileServer: fileServer,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, h.uriPrefix) // cut down the prefix for the file server
	h.fileServer.ServeHTTP(w, r)
	r.URL.Path = h.uriPrefix + r.URL.Path // restore the full path
}
