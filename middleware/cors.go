package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/raystack/raccoon/config"
)

var Cors func(http.Handler) http.Handler

func loadCors() {
	if config.ServerCors.Enabled {
		opts := []handlers.CORSOption{handlers.AllowedHeaders(config.ServerCors.AllowedHeaders),
			handlers.AllowedMethods(config.ServerCors.AllowedMethods),
			handlers.AllowedOrigins(config.ServerCors.AllowedOrigin)}
		if config.ServerCors.AllowCredentials {
			opts = append(opts, handlers.AllowCredentials())
		}
		if config.ServerCors.MaxAge > 0 {
			opts = append(opts, handlers.MaxAge(config.ServerCors.MaxAge))
		}
		Cors = handlers.CORS(opts...)
	} else {
		Cors = func(h http.Handler) http.Handler { return h }
	}

}
