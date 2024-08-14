package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/metrics"
	"github.com/raystack/raccoon/middleware"
	"github.com/raystack/raccoon/services/rest/websocket"
	"github.com/raystack/raccoon/services/rest/websocket/connection"
)

type Service struct {
	Collector collector.Collector
	s         *http.Server
}

func NewRestService(c collector.Collector) *Service {
	pingInterval := time.Duration(config.ServerWs.PingIntervalMS) * time.Millisecond
	writeWaitInterval := time.Duration(config.ServerWs.WriteWaitIntervalMS) * time.Millisecond
	pingChannel := make(chan connection.Conn, config.ServerWs.ServerMaxConn)
	wh := websocket.NewHandler(pingChannel, c)
	go websocket.Pinger(pingChannel, config.ServerWs.PingerSize, pingInterval, writeWaitInterval)

	go reportConnectionMetrics(*wh.Table())

	go websocket.AckHandler(websocket.AckChan)

	restHandler := NewHandler(c)
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(pingHandler).Methods(http.MethodGet)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", wh.HandlerWSEvents).Methods(http.MethodGet).Name("events")
	subRouter.HandleFunc("/events", restHandler.RESTAPIHandler).Methods(http.MethodPost).Name("events")

	server := &http.Server{
		Handler: applyMiddleware(router),
		Addr:    ":" + config.ServerWs.AppPort,
	}
	return &Service{
		s:         server,
		Collector: c,
	}
}

func applyMiddleware(router http.Handler) http.Handler {
	h := middleware.GetCors()(router)
	return h
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func reportConnectionMetrics(conn connection.Table) {
	interval := time.Duration(config.MetricInfo.RuntimeStatsRecordIntervalMS) * time.Millisecond
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
