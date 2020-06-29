package websocket

import (
	"context"
	"net/http"
	"raccoon/config"
	"raccoon/logger"

	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	HttpServer    *http.Server
	bufferChannel chan []*de.CSEventMessage
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
}

//CreateServer - instantiates the http server
func CreateServer() (*Server, chan []*de.CSEventMessage) {
	//create the websocket handler that upgrades the http request
	bufferChannel := make(chan []*de.CSEventMessage, config.WorkerConfigLoader().ChannelSize())
	wsHandler := &Handler{
		websocketUpgrader: getWebSocketUpgrader(config.ServerConfig.ReadBufferSize, config.ServerConfig.WriteBufferSize, config.ServerConfig.CheckOrigin),
		bufferChannel:     bufferChannel,
		user:              NewUserStore(config.ServerConfig.ServerMaxConn),
		PingInterval:      config.ServerConfig.PingInterval,
		PongWaitInterval:  config.ServerConfig.PongWaitInterval,
		WriteWaitInterval: config.ServerConfig.WriteWaitInterval,
	}
	server := &Server{
		HttpServer: &http.Server{
			Handler: Router(wsHandler),
			Addr:    ":" + config.ServerConfig.AppPort,
		},
		bufferChannel: bufferChannel,
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
