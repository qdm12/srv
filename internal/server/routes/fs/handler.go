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
	hasPrefix := strings.HasPrefix(r.URL.Path, h.uriPrefix)
	if hasPrefix {
		// cut down the prefix for the file server
		r.URL.Path = r.URL.Path[len(h.uriPrefix):]
	}

	h.fileServer.ServeHTTP(w, r)

	if hasPrefix {
		r.URL.Path = h.uriPrefix + r.URL.Path // restore the full path
	}
}
