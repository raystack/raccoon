# Raccoon Client Library for Go #

[![Go Reference](https://pkg.go.dev/badge/google.golang.org/api.svg)](https://pkg.go.dev/google.golang.org/api)

## Requirements
go 1.16 or above

## Install
```shell
go get github.com/odpf/raccoon/clients/go
```
## Usage

#### Construct a new REST client, then use the various options on the client.
For example:
```go
	client, err := rest.NewRest(
		rest.WithUrl("http://localhost:8080/api/v1/events"),
		rest.WithHeader("x-user-id", "123"),
		rest.WithSerializer(serializer.PROTO), // default is JSON
	)
```
#### Construct a new GRPC client, then use the various options on the client.
For example:
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