package pprof

import (
	"context"
	"net/http"
	// enable pprof https://pkg.go.dev/net/http/pprof#pkg-overview
	_ "net/http/pprof"
)

type Service struct {
	s *http.Server
}

func (s Service) Init() error {
	server := &http.Server{Addr: "localhost:6060", Handler: nil}
	s.s = server
	return server.ListenAndServe()
}

func (s Service) Name() string {
	return "pprof"
}

func (s Service) Shutdown(ctx context.Context) {
	s.s.Shutdown(ctx)
}
