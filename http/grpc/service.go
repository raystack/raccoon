package grpc

import (
	"context"
	"fmt"
	"net"
	"raccoon/collection"
	"raccoon/config"
	pb "raccoon/proto"

	"google.golang.org/grpc"
)

type Service struct {
	Collector collection.Collector
	s         *grpc.Server
}

func (s Service) Init(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.ServerGRPC.Port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterEventServiceServer(server, &Handler{C: s.Collector})
	s.s = server
	return server.Serve(lis)
}

func (s Service) Name() string {
	return "GRPC"
}

func (s Service) Shutdown(ctx context.Context) error {
	s.s.GracefulStop()
	return nil
}
