package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/odpf/raccoon/collection"
	"github.com/odpf/raccoon/config"
	pb "github.com/odpf/raccoon/proto"
	"google.golang.org/grpc"
)

type Service struct {
	Collector collection.Collector
	s         *grpc.Server
}

func NewGRPCService(c collection.Collector) *Service {
	server := grpc.NewServer()
	pb.RegisterEventServiceServer(server, &Handler{C: c})
	return &Service{
		s:         server,
		Collector: c,
	}
}

func (s *Service) Init(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.ServerGRPC.Port))
	if err != nil {
		return err
	}
	return s.s.Serve(lis)
}

func (s *Service) Name() string {
	return "GRPC"
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.s.GracefulStop()
	return nil
}
