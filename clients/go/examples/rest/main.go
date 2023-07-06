package main

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	raccoon "github.com/raystack/raccoon/clients/go"
	"github.com/raystack/raccoon/clients/go/rest"
	"github.com/raystack/raccoon/clients/go/serializer"
	"github.com/raystack/raccoon/clients/go/testdata"
)

func main() {

	client, err := rest.New(
		rest.WithUrl("http://localhost:8080/api/v1/events"),
		rest.WithHeader("x-user-id", "123"),
		rest.WithSerializer(serializer.PROTO), // default is JSON
		rest.WithRetry(time.Second*2, 5),
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
