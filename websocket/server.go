package websocket

import (
	"context"
	"net/http"
	"raccoon/config"
	"raccoon/logger"
	"runtime"
	"time"

	"raccoon/metrics"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "net/http/pprof"
)

type Server struct {
	HttpServer    *http.Server
	bufferChannel chan EventsBatch
	user          *User
	pingChannel   chan connection
}

func (s *Server) StartHTTPServer(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		logger.Info("WebSocket Server --> startHttpServer")
		err := s.HttpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Errorf("WebSocket Server --> HTTP Server could not be started = %s", err.Error())
			cancel()
		}
	}()
	go s.ReportServerMetrics()
	go Pinger(s.pingChannel, config.ServerConfig.PingerSize, config.ServerConfig.PingInterval, config.ServerConfig.WriteWaitInterval)
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
	t := time.Tick(config.StatsdConfigLoader().FlushPeriod())
	m := &runtime.MemStats{}
	for {
		<-t
		metrics.Gauge("connections.count", s.user.TotalUsers(), "")
		metrics.Gauge("go.routines.count", runtime.NumGoroutine(), "")

		runtime.ReadMemStats(m)
		metrics.Gauge("heap.alloc", m.HeapAlloc, "")
		metrics.Gauge("heap.inuse", m.HeapInuse, "")
		metrics.Gauge("heap.object", m.HeapObjects, "")
		metrics.Gauge("stack.inuse", m.StackInuse, "")
		metrics.Gauge("gc.triggered", m.LastGC/1000, "")
		metrics.Gauge("gc.pauseNs", m.PauseNs[(m.NumGC+255)%256]/1000, "")
		metrics.Gauge("gc.count", m.NumGC, "")
		metrics.Gauge("gc.pauseTotalNs", m.PauseTotalNs, "")
	}
}

//CreateServer - instantiates the http server
func CreateServer() (*Server, chan EventsBatch) {
	//create the websocket handler that upgrades the http request
	bufferChannel := make(chan EventsBatch, config.WorkerConfigLoader().ChannelSize())
	pingChannel := make(chan connection, config.ServerConfig.ServerMaxConn)
	user := NewUserStore(config.ServerConfig.ServerMaxConn)
	wsHandler := &Handler{
		websocketUpgrader: getWebSocketUpgrader(config.ServerConfig.ReadBufferSize, config.ServerConfig.WriteBufferSize, config.ServerConfig.CheckOrigin),
		bufferChannel:     bufferChannel,
		user:              user,
		PongWaitInterval:  config.ServerConfig.PongWaitInterval,
		WriteWaitInterval: config.ServerConfig.WriteWaitInterval,
		PingChannel:       pingChannel,
	}
	server := &Server{
		HttpServer: &http.Server{
			Handler: Router(wsHandler),
			Addr:    ":" + config.ServerConfig.AppPort,
		},
		bufferChannel: bufferChannel,
		user:          user,
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

func getWebSocketUpgrader(readBufferSize int, writeBufferSize int, checkOrigin bool) websocket.Upgrader {
	ug := websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return checkOrigin
		},
	}
	return ug
}
