package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/raystack/raccoon/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	url    = "ws://localhost:8080/api/v1/events"
	header = http.Header{
		"X-User-ID": []string{"1234"},
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
	request := &pb.SendEventRequest{
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
	SendEventResponse := &pb.SendEventResponse{}
	proto.Unmarshal(response, SendEventResponse)
	// Handle the response accordingly
	fmt.Printf("%v", SendEventResponse)
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
