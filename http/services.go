package http

import (
	"context"
	"raccoon/collection"
	"raccoon/logger"

	"raccoon/http/grpc"
	"raccoon/http/pprof"
	"raccoon/http/rest"
)

type bootstrapper interface {
	// Init initialize each HTTP based server. Return error if initialization failed. Put the Serve() function as return mostly suffice for Init process.
	Init(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Name() string
}

type Services struct {
	b []bootstrapper
}

func (s *Services) Start(ctx context.Context, cancel context.CancelFunc) {
	for _, init := range s.b {
		i := init
		go func() {
			logger.Infof("%s Server --> startServers", i.Name())
			err := i.Init(ctx)
			if err != nil {
				logger.Errorf("%s Server --> could not be started = %s", i.Name(), err)
				cancel()
			}
		}()
	}
}

func (s *Services) Shutdown(ctx context.Context) {
	for _, b := range s.b {
		logger.Infof("%s Server --> shutting down", b.Name())
		b.Shutdown(ctx)
	}
}

func Create(b chan collection.CollectRequest) Services {
	c := collection.NewChannelCollector(b)
	return Services{
		b: []bootstrapper{
			grpc.Service{Collector: c},
			pprof.Service{},
			rest.Service{Collector: c},
		},
	}
}
