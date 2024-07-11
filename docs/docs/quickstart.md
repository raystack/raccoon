import Tabs from '@theme/Tabs'
import TabItem from '@theme/TabItem'
import CodeBlock from '@theme/CodeBlock'

# Quickstart

This document will guide you on how to get Raccoon along with Kafka setup running locally. This document assumes that you have installed Docker and Kafka with `host.docker.internal` [advertised](https://www.confluent.io/blog/kafka-listeners-explained/) on your machine.

## Run Raccoon with Docker

Make sure to set `PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS` according to your local Kafka setup. Then run the following commands:. 

```bash
$ docker run -p 8080:8080 \
  -e SERVER_WEBSOCKET_CONN_ID_HEADER=X-User-ID \
  -e PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS=host.docker.internal:9092 \
  -e EVENT_DISTRIBUTION_PUBLISHER_PATTERN=clickstream-log \
  raystack/raccoon:latest
```

To test whether the service is running or not, you can try to ping the server.

```bash
$ curl http://localhost:8080/ping
```

## Publishing Your First Event

```mdx-code-block
<Tabs default>
<TabItem value='go'>
```
create a directory called `go-raccoon-example` and initalise a go module

```bash
$ mkdir go-raccoon-example
$ cd go-raccoon-example
$ go mod init go-raccoon-example
```
Install the raccoon client

``` bash
$ go get github.com/raystack/raccoon/clients/go
```
create the `main.go` file 
```go title="main.go" showLineNumbers
package main

import (
	"fmt"
	"log"
	"time"

	"crypto/tls"

	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	raccoon "github.com/raystack/raccoon/clients/go"
	"github.com/raystack/raccoon/clients/go/grpc"
	"github.com/raystack/raccoon/clients/go/testdata"
)

func main() {
	client, err := grpc.New(
		grpc.WithAddr("localhost:8080"),
		grpc.WithHeader("x-user-id", "123"),
		grpc.WithDialOptions(
			g.WithTransportCredentials(credentials.NewServerTLSFromCert(&tls.Certificate{})),
		), 
		grpc.WithRetry(time.Second*2, 5),
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
```

Finally, run the program
```bash
$ go run main.go
```

```mdx-code-block
</TabItem>
<TabItem value='cli'>CLI</TabItem>
<TabItem value='java'>Java</TabItem>
<TabItem value='js'>Javascript</TabItem>
</Tabs>
```


To verify the event published by Raccoon. First, you need to start a Kafka listener.

```bash
$ kafka-console-consumer --bootstrap-server localhost:9092 --topic clickstream-log
```

## Where To Go Next

For more detail about publishing events to Raccoon, you can read the [detailed document](guides/publishing.md) under the guides section. To understand more about how Raccoon works, you can go to the [architecture document](concepts/architecture.md).