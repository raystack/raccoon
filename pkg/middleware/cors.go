package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/raystack/raccoon/config"
)

var cors func(http.Handler) http.Handler

func loadCors() {
	if config.Server.CORS.Enabled {
		opts := []handlers.CORSOption{handlers.AllowedHeaders(config.Server.CORS.AllowedHeaders),
			handlers.AllowedMethods(config.Server.CORS.AllowedMethods),
			handlers.AllowedOrigins(config.Server.CORS.AllowedOrigin)}
		if config.Server.CORS.AllowCredentials {
			opts = append(opts, handlers.AllowCredentials())
		}
		if config.Server.CORS.MaxAge > 0 {
			opts = append(opts, handlers.MaxAge(config.Server.CORS.MaxAge))
		}
		cors = handlers.CORS(opts...)
	} else {
		cors = func(h http.Handler) http.Handler { return h }
	}
}

func GetCors() func(http.Handler) http.Handler {
	return cors
}
