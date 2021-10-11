package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"raccoon/config"
	raccoongrpc "raccoon/http/grpc"
	"raccoon/http/rest"
	"raccoon/http/websocket"
	"raccoon/http/websocket/connection"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/pkg/collection"
	pb "raccoon/pkg/proto"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type Servers struct {
	HTTPServer  *http.Server
	table       *connection.Table
	pingChannel chan connection.Conn
	GRPCServer  *grpc.Server
}

func (s *Servers) StartServers(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		logger.Info("HTTP Server --> startServers")
		err := s.HTTPServer.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Errorf("HTTP Server --> HTTP Server could not be started = %s", err.Error())
			cancel()
		}
	}()
	go func() {
		lis, err := net.Listen("tcp", config.ServerGRPC.Port)
		if err != nil {
			logger.Errorf("GRPC Server --> GRPC Server could not be started = %s", err.Error())
			cancel()
		}
		logger.Info("GRPC Server -> startServers")
		if err := s.GRPCServer.Serve(lis); err != nil {
			logger.Errorf("GRPC Server --> GRPC Server could not be started = %s", err.Error())
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

func (s *Servers) ReportServerMetrics() {
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
func CreateServer(bufferChannel chan *collection.EventsBatch) *Servers {
	//create the websocket handler that upgrades the http request
	collector := collection.NewChannelCollector(bufferChannel)
	pingChannel := make(chan connection.Conn, config.ServerWs.ServerMaxConn)
	wsHandler := websocket.NewHandler(pingChannel)
	restHandler := &rest.Handler{}
	grpcHandler := &raccoongrpc.Handler{}
	handler := &Handler{wsHandler, restHandler, grpcHandler}
	grpcServer := grpc.NewServer(grpc.WithInsecure())
	server := &Servers{
		HTTPServer: &http.Server{
			Handler: Router(handler, collector),
			Addr:    ":" + config.ServerWs.AppPort,
		},
		table:       wsHandler.Table(),
		pingChannel: pingChannel,
		GRPCServer:  grpcServer,
	}
	pb.RegisterEventServiceServer(grpcServer, handler.gh)
	//Wrap the handler with a Server instance and return it
	return server
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// Router sets up the routes
func Router(h *Handler, collector collection.Collector) http.Handler {
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(PingHandler).Methods(http.MethodGet)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", h.wh.GetHandlerWSEvents(collector)).Methods(http.MethodGet).Name("events")
	subRouter.HandleFunc("/events", h.rh.GetRESTAPIHandler(collector)).Methods(http.MethodPost).Name("events")
	return router
}
