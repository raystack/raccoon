package websocket

import (
	"clickstream-service/config"
	"clickstream-service/logger"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	server *negroni.Negroni
}

func (s *Server) StartHTTPServer(ctx context.Context, cancel context.CancelFunc) {
	port := fmt.Sprintf(":%s",config.AppPort())
	go s.server.Run(port)
	logger.Info("WebSocket Server --> startHttpServer")
	go shutDownGracefully(ctx, cancel)
}

//CreateServer - instantiates the http server
func CreateServer() *Server {
	//create the websocket handler that upgrades the http request
	wsHandler := &Handler{
		websocketUpgrader: getWebSocketUpgrader(),
		//@TODO - init with the events-channel
	}
	negRoniServer := negroni.New(negroni.NewRecovery())
	//create & set the router
	negRoniServer.UseHandler(Router(wsHandler))
	//Wrap the handler with a Server instance and return it
	return &Server{
		server: negRoniServer,
	}

}

// Router sets up the routes
func Router(h *Handler) http.Handler {
	router := mux.NewRouter()
	router.Path("/ping").HandlerFunc(PingHandler).Methods(http.MethodGet)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", h.HandlerWSEvents).Methods(http.MethodGet).Name("events")
	return router
}

func getWebSocketUpgrader() websocket.Upgrader {
	/**
	@TODO - should make the buffer sizes & cross-origin configurable
	*/
	ug := websocket.Upgrader{
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return ug
}

func shutDownGracefully(ctx context.Context, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[Websocket.server] Received signal %s, shutting down http server", sig))
			//@TODO - Should see a way to stop the http negroni server
		default:
			logger.Info(fmt.Sprintf("[Websocket.server] Received a unexpected signal %s", sig))
		}
	}
}
