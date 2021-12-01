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
	Init() error
	Shutdown(ctx context.Context)
	Name() string
}

type Servers struct {
	b []bootstrapper
}

func (s *Servers) StartServers(ctx context.Context, cancel context.CancelFunc) {
	for _, init := range s.b {
		i := init
		go func() {
			logger.Infof("%s Server --> startServers", i.Name())
			err := i.Init()
			if err != nil {
				logger.Errorf("%s Server --> could not be started = %s", i.Name(), err)
				cancel()
			}
		}()
	}
}

func (s *Servers) Shutdown(ctx context.Context) {
	for _, b := range s.b {
		logger.Infof("%s Server --> shutting down", b.Name())
		b.Shutdown(ctx)
	}
}

func CreateServer(b chan *collection.CollectRequest) Servers {
	return Servers{
		b: []bootstrapper{
			grpc.Service{Buffer: b},
			pprof.Service{},
			rest.Service{Buffer: b},
		},
	}
}
