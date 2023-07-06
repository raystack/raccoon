package main

import (
	"fmt"
	"log"

	raccoon "github.com/raystack/raccoon/clients/go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/raystack/raccoon/clients/go/testdata"
	"github.com/raystack/raccoon/clients/go/ws"
)

func main() {
	client, err := ws.New(
		ws.WithUrl("ws://localhost:8080/api/v1/events"),
		ws.WithHeader("x-user-id", "123"),
		ws.WithHeader("x-user-type", "gojek"))

	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	reqGuid, err := client.Send([]*raccoon.Event{
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
	resp := <-client.EventAcks()
	fmt.Println(resp.Status)
	fmt.Println(resp.Data["req_guid"])
}
