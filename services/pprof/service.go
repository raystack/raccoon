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

func NewPprofService() *Service {
	return &Service{
		s: &http.Server{Addr: "localhost:6060", Handler: nil},
	}
}

func (s *Service) Init(ctx context.Context) error {
	return s.s.ListenAndServe()
}

func (s *Service) Name() string {
	return "pprof"
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
