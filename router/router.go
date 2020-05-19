package router

import (
	"net/http"

	"clickstream-service/handler"
	"github.com/gorilla/mux"
)

func Router() http.Handler {
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(handler.PingHandler).Methods(http.MethodGet)

	//subRouter := router.PathPrefix("/api/v1").Subrouter()
	//subRouter.HandleFunc("/events", handler.CreateEvent()).Methods(http.MethodPost).Name("create-event")

	return router
}
