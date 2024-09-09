package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/core/collector"
	"github.com/raystack/raccoon/pkg/metrics"
	"github.com/raystack/raccoon/server/rest/websocket"
	"github.com/raystack/raccoon/server/rest/websocket/connection"
)

type Service struct {
	Collector collector.Collector
	s         *http.Server
}

func NewRestService(c collector.Collector) *Service {
	pingInterval := time.Duration(config.Server.Websocket.PingIntervalMS) * time.Millisecond
	writeWaitInterval := time.Duration(config.Server.Websocket.WriteWaitIntervalMS) * time.Millisecond
	pingChannel := make(chan connection.Conn, config.Server.Websocket.ServerMaxConn)
	wh := websocket.NewHandler(pingChannel, c)
	go websocket.Pinger(pingChannel, config.Server.Websocket.PingerSize, pingInterval, writeWaitInterval)

	go reportConnectionMetrics(*wh.Table())

	go websocket.AckHandler(websocket.AckChan)

	restHandler := NewHandler(c)
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(pingHandler).Methods(http.MethodGet)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", wh.HandlerWSEvents).Methods(http.MethodGet).Name("events")
	subRouter.HandleFunc("/events", restHandler.RESTAPIHandler).Methods(http.MethodPost).Name("events")

	server := &http.Server{
		Handler: withCORS(router),
		Addr:    ":" + config.Server.Websocket.AppPort,
	}
	return &Service{
		s:         server,
		Collector: c,
	}
}

func withCORS(router http.Handler) http.Handler {
	if !config.Server.CORS.Enabled {
		return router
	}
	opts := []handlers.CORSOption{handlers.AllowedHeaders(config.Server.CORS.AllowedHeaders),
		handlers.AllowedMethods(config.Server.CORS.AllowedMethods),
		handlers.AllowedOrigins(config.Server.CORS.AllowedOrigin)}
	if config.Server.CORS.AllowCredentials {
		opts = append(opts, handlers.AllowCredentials())
	}
	if config.Server.CORS.MaxAge > 0 {
		opts = append(opts, handlers.MaxAge(config.Server.CORS.MaxAge))
	}
	return handlers.CORS(opts...)(router)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func reportConnectionMetrics(conn connection.Table) {
	interval := time.Duration(config.Metric.RuntimeStatsRecordIntervalMS) * time.Millisecond
	for range time.Tick(interval) {
		for k, v := range conn.TotalConnectionPerGroup() {
			metrics.Gauge("connections_count_current", v, map[string]string{"conn_group": k})
		}
	}
}

func (s *Service) Init(context.Context) error {
	return s.s.ListenAndServe()
}

func (*Service) Name() string {
	return "REST"
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
