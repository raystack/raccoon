package main

import (
	"fmt"
	"log"

	"crypto/tls"

	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/odpf/raccoon/clients/go/grpc"
	"github.com/odpf/raccoon/clients/go/testdata"
)

func main() {
	client, err := grpc.NewGrpc(
		grpc.WithAddr("http://localhost:8080"),
		grpc.WithHeader("x-user-id", "123"),

		grpc.WithDialOptions(
			g.WithTransportCredentials(credentials.NewServerTLSFromCert(&tls.Certificate{})),
		), // default is insecure

		// default serializer is proto.
	)

	if err != nil {
		log.Fatal(err)
	}

	reqGuid, resp, err := client.Send([]*raccoon.Event{
		{
			Type: "page",
			Data: &testdata.PageEvent{
				EventGuid: uuid.NewString(),
				EventName: "clicked",
				SentTime:  timestamppb.Now(),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reqGuid)
	fmt.Println(resp.Status)
}
