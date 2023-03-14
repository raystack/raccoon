package grpc

import (
	"context"
	"fmt"
	"net"

	pbgrpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/raccoon/v1beta1/raccoonv1beta1grpc"
	"github.com/goto/raccoon/collection"
	"github.com/goto/raccoon/config"
	"google.golang.org/grpc"
)

type Service struct {
	Collector collection.Collector
	s         *grpc.Server
}

func NewGRPCService(c collection.Collector) *Service {
	server := grpc.NewServer()
	pbgrpc.RegisterEventServiceServer(server, &Handler{C: c})
	return &Service{
		s:         server,
		Collector: c,
	}
}

func (s *Service) Init(context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.ServerGRPC.Port))
	if err != nil {
		return err
	}
	return s.s.Serve(lis)
}

func (*Service) Name() string {
	return "GRPC"
}

func (s *Service) Shutdown(context.Context) error {
	s.s.GracefulStop()
	return nil
}
