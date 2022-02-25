package rest

import (
	"context"
	"fmt"
	"net/http"
	"raccoon/collection"
	"raccoon/config"
	"raccoon/services/rest/websocket"
	"raccoon/services/rest/websocket/connection"
	"raccoon/metrics"
	"time"

	"github.com/gorilla/mux"
)

type Service struct {
	Collector collection.Collector
	s         *http.Server
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func reportConnectionMetrics(conn connection.Table) {
	t := time.Tick(config.MetricStatsd.FlushPeriodMs)
	for {
		<-t
		for k, v := range conn.TotalConnectionPerGroup() {
			metrics.Gauge("connections_count_current", v, fmt.Sprintf("conn_group=%s", k))
		}
	}
}

func (s Service) Init(ctx context.Context) error {
	pingChannel := make(chan connection.Conn, config.ServerWs.ServerMaxConn)
	wh := websocket.NewHandler(pingChannel, s.Collector)
	go websocket.Pinger(pingChannel, config.ServerWs.PingerSize, config.ServerWs.PingInterval, config.ServerWs.WriteWaitInterval)

	go reportConnectionMetrics(*wh.Table())

	restHandler := NewHandler(s.Collector)
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(pingHandler).Methods(http.MethodGet)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", wh.HandlerWSEvents).Methods(http.MethodGet).Name("events")
	subRouter.HandleFunc("/events", restHandler.RESTAPIHandler).Methods(http.MethodPost).Name("events")

	server := &http.Server{
		Handler: router,
		Addr:    ":" + config.ServerWs.AppPort,
	}
	s.s = server
	return server.ListenAndServe()
}

func (s Service) Name() string {
	return "REST"
}

func (s Service) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
