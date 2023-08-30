package middleware

import "net/http"

var loaded bool

func Load() {
	if loaded {
		return
	}
	loadCors()
	loaded = true
}

func Apply(h http.Handler) http.Handler {
	handler := cors(h)
	return handler
}
