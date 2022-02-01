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
	Buffer chan collection.CollectRequest
	s      *grpc.Server
}

func (s Service) Init() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.ServerGRPC.Port))
	if err != nil {
		return err
	}
	collector := collection.NewChannelCollector(s.Buffer)
	server := grpc.NewServer()
	pb.RegisterEventServiceServer(server, &Handler{C: collector})
	s.s = server
	return server.Serve(lis)
}

func (s Service) Name() string {
	return "GRPC"
}

func (s Service) Shutdown(ctx context.Context) {
	s.s.GracefulStop()
}
