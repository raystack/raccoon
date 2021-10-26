package websocket

import (
	"context"
	"fmt"
	"net/http"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/websocket/connection"
	"runtime"
	"time"

	"raccoon/metrics"

	"github.com/gorilla/mux"

	// https://golang.org/pkg/net/http/pprof/
	_ "net/http/pprof"
)

type Server struct {
	HTTPServer    *http.Server
	bufferChannel chan EventsBatch
	table         *connection.Table
	pingChannel   chan connection.Conn
}

func (s *Server) StartHTTPServer(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		logger.Info("WebSocket Server --> startHttpServer")
		err := s.HTTPServer.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Errorf("WebSocket Server --> HTTP Server could not be started = %s", err.Error())
			cancel()
		}
	}()
	go s.ReportServerMetrics()
	go Pinger(s.pingChannel, config.ServerWs.PingerSize, config.ServerWs.PingInterval, config.ServerWs.WriteWaitInterval)
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logger.Errorf("WebSocket Server --> pprof could not be enabled: %s", err.Error())
			cancel()
		} else {
			logger.Info("WebSocket Server --> pprof :: Enabled")
		}
	}()
}

func (s *Server) ReportServerMetrics() {
	t := time.Tick(config.MetricStatsd.FlushPeriodMs)
	m := &runtime.MemStats{}
	for {
		<-t
		for k, v := range s.table.TotalConnectionPerGroup() {
			metrics.Gauge("connections_count_current", v, fmt.Sprintf("conn_group=%s", k))
		}
		metrics.Gauge("server_go_routines_count_current", runtime.NumGoroutine(), "")

		runtime.ReadMemStats(m)
		metrics.Gauge("server_mem_heap_alloc_bytes_current", m.HeapAlloc, "")
		metrics.Gauge("server_mem_heap_inuse_bytes_current", m.HeapInuse, "")
		metrics.Gauge("server_mem_heap_objects_total_current", m.HeapObjects, "")
		metrics.Gauge("server_mem_stack_inuse_bytes_current", m.StackInuse, "")
		metrics.Gauge("server_mem_gc_triggered_current", m.LastGC/1000, "")
		metrics.Gauge("server_mem_gc_pauseNs_current", m.PauseNs[(m.NumGC+255)%256]/1000, "")
		metrics.Gauge("server_mem_gc_count_current", m.NumGC, "")
		metrics.Gauge("server_mem_gc_pauseTotalNs_current", m.PauseTotalNs, "")
	}
}

//CreateServer - instantiates the http server
func CreateServer() (*Server, chan EventsBatch) {
	//create the websocket handler that upgrades the http request
	bufferChannel := make(chan EventsBatch, config.Worker.ChannelSize)
	pingChannel := make(chan connection.Conn, config.ServerWs.ServerMaxConn)
	ugConfig := connection.UpgraderConfig{
		ReadBufferSize:    config.ServerWs.ReadBufferSize,
		WriteBufferSize:   config.ServerWs.WriteBufferSize,
		CheckOrigin:       config.ServerWs.CheckOrigin,
		MaxUser:           config.ServerWs.ServerMaxConn,
		PongWaitInterval:  config.ServerWs.PongWaitInterval,
		WriteWaitInterval: config.ServerWs.WriteWaitInterval,
		ConnIDHeader:      config.ServerWs.ConnIDHeader,
		ConnGroupHeader:   config.ServerWs.ConnGroupHeader,
		ConnGroupDefault:  config.ServerWs.ConnGroupDefault,
	}
	upgrader := connection.NewUpgrader(ugConfig)
	wsHandler := &Handler{
		upgrader:      upgrader,
		bufferChannel: bufferChannel,
		PingChannel:   pingChannel,
	}
	server := &Server{
		HTTPServer: &http.Server{
			Handler: Router(wsHandler),
			Addr:    ":" + config.ServerWs.AppPort,
		},
		table:         upgrader.Table,
		bufferChannel: bufferChannel,
		pingChannel:   pingChannel,
	}
	//Wrap the handler with a Server instance and return it
	return server, bufferChannel
}

// Router sets up the routes
func Router(h *Handler) http.Handler {
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(PingHandler).Methods(http.MethodGet)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", h.HandlerWSEvents).Methods(http.MethodGet).Name("events")
	return router
}
