// Package fs is the HTTP handler to serve files from an http.FileSystem.
package fs

import (
	"net/http"
)

func NewHandler(fs http.FileSystem) http.Handler {
	return http.FileServer(fs)
}
