# Golang

## Requirements
Make sure that Go >= `1.16` is installed on your system. See [installation instructions](https://go.dev/doc/install) on Go's website for more info.

## Installation
Install Raccoon's Go client using
```bash
$ go get github.com/raystack/raccoon/clients/go
```
## Usage

### Quickstart

Below is a self contained example of Raccoon's Go client that uses the Websocket API to publish events

```go title="quickstart.go"
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
        ws.WithHeader("x-user-type", "ACME"))
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
Depending on which protocol you wish to utilise to publish events to Raccoon, you will need to intantiate a different client.

Following is a table describing which package you should use for a given protocol.

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

#### Sending events

Event's can be sent using `client.Send(events []*raccoon.Event)`. The return signature of the `Send` method depends on the type of Client.

| Type | Signature |
| --- | --- |
| `REST`, `gRPC` | `Send([]*raccoon.Event) (string, *raccoon.Response, error)` |
| `Websocket` | `Send([]*raccoon.Event) (string, error)` |

For `gRPC` and `REST` clients, the response is returned synchronously. For `Websocket` the responses are returned asynchronously via a channel returned by `EventAcks()`.

### Examples
You can find examples of clients over different protocols [here](https://github.com/raystack/raccoon/tree/main/clients/go/examples)