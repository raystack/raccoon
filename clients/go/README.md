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
client, err := rest.New(
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
client, err := grpc.New(
	grpc.WithAddr("http://localhost:8080"),
	grpc.WithHeader("x-user-id", "123"),
	grpc.WithDialOptions(
		g.WithTransportCredentials(credentials.NewServerTLSFromCert(&tls.Certificate{})),
	), // default is insecure
	// default serializer is proto.
)
```

see example: [examples/grpc](examples/grpc/main.go)

#### Construct a new Websocket client, then use the various options on the client.
For example:
```go
import "github.com/odpf/raccoon/clients/go/ws"
```
```go
client, err := ws.New(
	ws.WithUrl("ws://localhost:8080/api/v1/events"),
	ws.WithHeader("x-user-id", "123"),
	ws.WithHeader("x-user-type", "gojek"))
	// default serializer is proto.
```
Reading the message acknowledgements
```go
resp := <-client.EventAcks()
```
see example: [examples/websocket](examples/ws/main.go)

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



#### Retry Configuration
Default settings, wait = `1 second`, and maximum attempts = `3`. The following options can override it.

```go
rest.WithRetry(time.Second*2, 5)

grpc.WithRetry(time.Second*2, 5)
```

### Custom Logging Configuration
The default logger logs the information to the standard output,
and the default logger can be diabled by the following settings
```go
rest.WithLogger(nil)
grpc.WithLogger(nil)
```

The client provide the logger interface that can be implemented by any logger.
```go
type Logger interface {
	Infof(msg string, keysAndValues ...interface{})
	Errorf(msg string, keysAndValues ...interface{})
}
```
And cutomer logger can set for the client with the following options.

```go
rest.WithLogger(&CustomLogger{})
grpc.WithLogger(&CustomLogger{})
```
