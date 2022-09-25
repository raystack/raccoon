package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/odpf/raccoon/clients/go/rest"
	"github.com/odpf/raccoon/clients/go/serializer"
	"github.com/odpf/raccoon/clients/go/testdata"
)

func main() {

	client, err := rest.NewRest(
		rest.WithUrl("http://localhost:8080/api/v1/events"),
		rest.WithHeader("x-user-id", "123"),
		rest.WithSerializer(serializer.PROTO), // default is JSON
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
