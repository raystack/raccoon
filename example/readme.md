There are couple of things the client to handle to start send events to Raccoon.

- [Establish Websocket Connection](#establish-websocket-connection)
- [Set Ping/Pong Handler](#set-pingpong-handler)
- [Batch Events](#batch-events)
- [Send The Batch](#send-the-batch)
- [Handle The Response](#handle-the-response)

Below are the explanation of sample client in [main.go](https://github.com/raystack/raccoon/tree/main/docs/example/main.go)

## Establish Websocket Connection

You are free to use any websocket client as long as it supports passing header. You can connect to `/api/v1/events` endpoint with uniq id header set. You'll also need to handle retry in case Raccon reject the connection because [max connection is reached]().

```go
var (
	url    = "ws://localhost:8080/api/v1/events"
	header = http.Header{
		"X-User-ID": []string{"1234"},
	}
)

func main() {
	ws, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		panic("Fail to make websocket connection")
	}
```

## Set Ping/Pong Handler

Raccoon needs to maintain many alive connections. To clean up dead connection, Raccoon will close client that couldn't response to Ping on time. Because of that, the client need to handle Ping if not handled by the websocket client.

```go
	// Gorilla websocket has default PingHandler which sends Pong. No need to explicitly heandle it.
	go pinger(ws)
```

You can also check the liveliness of the server by having Pinger function and close the connection if necessary

```go
func pinger(ws *websocket.Conn) {
	ticker := time.Tick(pingInterval)
	for {
		<-ticker
		ws.WriteControl(websocket.PingMessage, []byte("--ping--"), time.Now().Add(pingInterval))
		fmt.Println("ping")
	}
}
```

## Batch Events

When the connection is set, all you need to do is collect the event and send them in batch.

```go
	event1 := generateSampleEvent()
	event2 := generateSampleEvent()
	eventBatch := []*pb.Event{
		event1,
		event2,
	}
```

Where `generateSampleEvent` is

```go
func generateSampleEvent() *pb.Event {
	sampleEvent := &SampleEvent{Description: "user_click"}
	sampleBin, _ := proto.Marshal(sampleEvent)
	event := &pb.Event{EventBytes: sampleBin, Type: "some-type"}
	return event
}
```

## Send The Batch

Now you have websocket connection and batch of event ready, all you need is send the batch. Don't forget to fill `send_time` field before sending the request.

```go
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
```

## Handle The Response

Raccoon sends SendEventResponse for every batch of events sent from the client. The ReqGuid in the response identifies the batch that the client sent. The response object could be used for client's telemetry in terms of how may batches succeeded, failed etc., The clients can retry based on failures however server side kafka send failures are not sent as failures due to the [acknowledgement design as explained here](https://github.com/raystack/raccoon/blob/main/docs/concepts/architecture.md#acknowledging-events).

```go
	_, response, _ := ws.ReadMessage()
	SendEventResponse := &pb.SendEventResponse{}
	proto.Unmarshal(response, SendEventResponse)
	// Handle the response accordingly
	fmt.Printf("%v", SendEventResponse)
```
