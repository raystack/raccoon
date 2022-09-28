# Raccoon Client Library for Go #

[![Go Reference](https://pkg.go.dev/badge/google.golang.org/api.svg)](https://pkg.go.dev/google.golang.org/api)

## Requirements
go 1.16 or above

## Install
```bash
go get github.com/odpf/raccoon/clients/go
```
## Usage

#### Construct a new REST client, then use the various options on the client.
For example:
```go
import "github.com/odpf/raccoon/clients/go/rest"
```
```go
	client, err := rest.NewRest(
		rest.WithUrl("http://localhost:8080/api/v1/events"),
		rest.WithHeader("x-user-id", "123"),
		rest.WithSerializer(serializer.PROTO), // default is JSON
	)
```

see example: [examples/rest](examples/rest/main.go)

#### Construct a new GRPC client, then use the various options on the client.
For example:
```go
import "github.com/odpf/raccoon/clients/go/grpc"
```
```go
	client, err := grpc.NewGrpc(
		grpc.WithAddr("http://localhost:8080"),
		grpc.WithHeader("x-user-id", "123"),

		grpc.WithDialOptions(
			g.WithTransportCredentials(credentials.NewServerTLSFromCert(&tls.Certificate{})),
		), // default is insecure

		// default serializer is proto.
	)
```

see example: [examples/grpc](examples/grpc/main.go)

#### Sending the request to raccoon
```go
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
```