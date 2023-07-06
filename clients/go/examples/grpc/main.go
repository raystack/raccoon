package main

import (
	"fmt"
	"log"
	"time"

	"crypto/tls"

	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	raccoon "github.com/raystack/raccoon/clients/go"
	"github.com/raystack/raccoon/clients/go/grpc"
	"github.com/raystack/raccoon/clients/go/testdata"
)

func main() {
	client, err := grpc.New(
		grpc.WithAddr("localhost:8080"),
		grpc.WithHeader("x-user-id", "123"),

		grpc.WithDialOptions(
			g.WithTransportCredentials(credentials.NewServerTLSFromCert(&tls.Certificate{})),
		), // default is insecure

		// default serializer is proto.

		grpc.WithRetry(time.Second*2, 5),
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
