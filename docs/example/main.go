package main

import (
	"fmt"
	"net/http"
	pb "raccoon/websocket/proto"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	url    = "ws://localhost:8080/api/v1/events"
	header = http.Header{
		"x-user-id": []string{"1234"},
	}
	pingInterval = 5 * time.Second
)

func main() {
	ws, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		panic("Fail to make websocket connection")
	}
	// Gorilla websocket has default PingHandler which sends Pong. No need to explicitly heandle it.
	go pinger(ws)

	event1 := generateSampleEvent()
	event2 := generateSampleEvent()
	eventBatch := []*pb.Event{
		event1,
		event2,
	}

	sentTime := time.Now()
	request := &pb.EventRequest{
		ReqGuid: "55F648D1-9A73-4F6C-8657-4D26A6C1F168",
		SentTime: &timestamppb.Timestamp{
			Seconds: sentTime.Unix(),
			Nanos:   int32(sentTime.Nanosecond()),
		},
		Events: eventBatch,
	}
	reqBin, _ := proto.Marshal(request)
	ws.WriteMessage(websocket.BinaryMessage, reqBin)

	_, response, _ := ws.ReadMessage()
	eventResponse := &pb.EventResponse{}
	proto.Unmarshal(response, eventResponse)
	// Handle the response accordingly
	fmt.Printf("%v", eventResponse)
}

func generateSampleEvent() *pb.Event {
	sampleEvent := &SampleEvent{Description: "user_click"}
	sampleBin, _ := proto.Marshal(sampleEvent)
	event := &pb.Event{EventBytes: sampleBin, Type: "some-type"}
	return event
}

func pinger(ws *websocket.Conn) {
	ticker := time.Tick(pingInterval)
	for {
		<-ticker
		ws.WriteControl(websocket.PingMessage, []byte("--ping--"), time.Now().Add(pingInterval))
		fmt.Println("ping")
	}
}
