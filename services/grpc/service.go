package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/config"
	pb "github.com/raystack/raccoon/proto"
	"google.golang.org/grpc"
)

type Service struct {
	Collector collector.Collector
	s         *grpc.Server
}

func NewGRPCService(c collector.Collector) *Service {
	server := grpc.NewServer()
	pb.RegisterEventServiceServer(server, &Handler{C: c})
	return &Service{
		s:         server,
		Collector: c,
	}
}

func (s *Service) Init(context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Server.GRPC.Port))
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
