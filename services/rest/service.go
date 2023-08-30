package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/raystack/raccoon/collection"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/metrics"
	"github.com/raystack/raccoon/middleware"
	"github.com/raystack/raccoon/services/rest/websocket"
	"github.com/raystack/raccoon/services/rest/websocket/connection"
)

type Service struct {
	Collector collection.Collector
	s         *http.Server
}

func NewRestService(c collection.Collector) *Service {
	pingChannel := make(chan connection.Conn, config.ServerWs.ServerMaxConn)
	wh := websocket.NewHandler(pingChannel, c)
	go websocket.Pinger(pingChannel, config.ServerWs.PingerSize, config.ServerWs.PingInterval, config.ServerWs.WriteWaitInterval)

	go reportConnectionMetrics(*wh.Table())

	go websocket.AckHandler(websocket.AckChan)

	restHandler := NewHandler(c)
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(pingHandler).Methods(http.MethodGet)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", wh.HandlerWSEvents).Methods(http.MethodGet).Name("events")
	subRouter.HandleFunc("/events", restHandler.RESTAPIHandler).Methods(http.MethodPost).Name("events")

	server := &http.Server{
		Handler: middleware.Apply(router),
		Addr:    ":" + config.ServerWs.AppPort,
	}
	return &Service{
		s:         server,
		Collector: c,
	}
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

func (s *Service) Init(context.Context) error {
	return s.s.ListenAndServe()
}

func (*Service) Name() string {
	return "REST"
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
