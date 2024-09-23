# Golang

## Requirements
Make sure that Go >= `1.16` is installed on your system. See [installation instructions](https://go.dev/doc/install) on Go's website for more info.

## Installation
Install Raccoon's Go client using [go get](https://go.dev/ref/mod#go-get)
```bash
$ go get github.com/raystack/raccoon/clients/go
```
## Usage

### Quickstart

Below is a self contained example of Raccoon's Go client that uses the Websocket API to publish events

```go title="quickstart.go" showLineNumbers
package main

import (
    "fmt"
    "log"
    raccoon "github.com/raystack/raccoon/clients/go"
    "google.golang.org/protobuf/types/known/timestamppb"
    "github.com/google/uuid"
    "github.com/raystack/raccoon/clients/go/serializer"
    "github.com/raystack/raccoon/clients/go/testdata"
    "github.com/raystack/raccoon/clients/go/ws"
)

func main() {
    client, err := ws.New(
        ws.WithUrl("ws://localhost:8080/api/v1/events"),
        ws.WithHeader("x-user-id", "123"),
        ws.WithHeader("x-user-type", "ACME"),
        ws.WithSerializer(serializer.JSON),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    _, err = client.Send([]*raccoon.Event{
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
    fmt.Println(<-client.EventAcks())
}
```

### Guide

#### Creating a client

Raccoon's API is exposed over 3 different protocols.
Depending on which protocol you wish to utilise to publish events to Raccoon, you will need to instantiate a different client.

Following is a table describing which client package you should use for a given protocol.

| Protocol | Package |
| --- | --- |
| Websocket | `github.com/raystack/raccoon/clients/go/ws` |
| REST | `github.com/raystack/raccoon/clients/go/rest` |
| gRPC | `github.com/raystack/raccoon/clients/go/grpc` |

For instance, you can create a client over REST API using:
```go
import "github.com/raystack/raccoon/clients/go/rest"

func main() {
    client, err := rest.New(
        rest.WithURL("http://localhost:8080/api/v1/events"),
        rest.WithHeader("x-user-id", "123")
    )
    if err != nil {
        panic(err)
    }

    // use the client here
}
```

Depending on which protocol client you create, specifying the URL or the address of the server is mandatory.

For `REST` and `Websocket` clients, this can be done via the `WithUrl` option. For `gRPC` server you must use the `WithAddr` option.

#### Sending events

Event's can be sent using `client.Send(events []*raccoon.Event)`. The return signature of the `Send` method depends on the type of Client.

| Type | Signature |
| --- | --- |
| `REST`, `gRPC` | `Send([]*raccoon.Event) (string, *raccoon.Response, error)` |
| `Websocket` | `Send([]*raccoon.Event) (string, error)` |

For `gRPC` and `REST` clients, the response is returned synchronously. For `Websocket` the responses are returned asynchronously via a channel returned by `EventAcks()`.

`Event` struct has two fields: `Type` and `Data`.
`Type` denotes the event type. This is used by raccoon to route the event to a specific topic downstream. `Data` field contains the payload. This data is serialised by the `serializer` that's configured on the client. The serializer can be configured by using the `WithSerializer()` option of the respective clients.

The following table lists which serializer to use for a given payload type.

| Message Type | Serializer |
| --- | --- |
| JSON | `serializer.JSON` |
| Protobuf | `serializer.PROTO`|

Once a client is constructed with a specific kind of serializer, you may only pass it events of that specific type. In particular, for `JSON` serialiser the event data must be a value that can be encoded by [`json.Marshal`](https://pkg.go.dev/encoding/json#Marshal). While for `PROTOBUF` serialiser the event data must be a protobuf message.

### Examples
You can find examples of client usage over different protocols [here](https://github.com/raystack/raccoon/tree/main/clients/go/examples)