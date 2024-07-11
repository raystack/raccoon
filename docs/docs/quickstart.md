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

To verify the event published by Raccoon. First, you need to start a Kafka listener. In a seperate terminal run:

```bash
$ kafka-console-consumer --bootstrap-server localhost:9092 --topic clickstream-log
```

## Publishing Your First Event

```mdx-code-block
<Tabs default>
<TabItem value='go'>
```
Make sure that `Go` >= `1.16` is installed on your system. See [installation instructions](https://go.dev/doc/install) on Go's website for more info.

Create a directory called `go-raccoon-example` and initalise it as a go module

```bash
$ mkdir go-raccoon-example
$ cd go-raccoon-example
$ go mod init go-raccoon-example
```
Install the raccoon client

``` bash
$ go get github.com/raystack/raccoon/clients/go
```
Create the `main.go` file 
```go title="main.go" showLineNumbers
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

	<-client.EventAcks()
}
```

Finally, run the program
```bash
$ go run main.go
```

```mdx-code-block
</TabItem>
<TabItem value='terminal'>
```
Make sure you have `curl` installed. On a debian-based system you can install `curl` via:

```bash
$ sudo apt-get install -y curl
```

Run the following to publish a single event to Raccoon
```bash
$ curl -XPOST "http://localhost:8080/api/v1/events" \
    -H "content-type: application/json" \
    -H "X-User-ID: user-one" \
    -d "
{
    \"req_guid\": \"foobar-123\",
    \"sent_time\": {
        \"seconds\": $(date +%s),
        \"nanos\": $(date +%N)
    },
    \"events\": [
        {
            \"type\": \"page\",
            \"eventBytes\": \"$(echo \"EVENT\" | base64)\"
        }
    ]
}"
```

```mdx-code-block
</TabItem>
<TabItem value='java'>Java</TabItem>
```

```mdx-code-block
<TabItem value='js'>
```
Make sure you have `node` >= `20.x` installed. See [installation instructions](https://nodejs.org/en/download/package-manager) on nodejs website for more info.

Create a new folder called `js-raccoon-example` and initalise it as a npm package.

```bash
$ mkdir js-raccoon-example
$ cd js-raccoon-example
$ npm init
```

Install the client using:
```bash
$ npm install @raystack/raccoon --save
```

Create a `main.js` file with the following contents:
```js title="main.js" showLineNumbers
import { RaccoonClient, SerializationType, WireType } from '@raystack/raccoon';

const logger = console;

//  create json messages
const jsonEvents = [
    {
        type: 'test-topic1',
        data: { key1: 'value1', key2: ['a', 'b'] }
    },
    {
        type: 'test-topic2',
        data: { key1: 'value2', key2: { key3: 'value3', key4: 'value4' } }
    }
];

//  initialise the raccoon client with required configs
const raccoonClient = new RaccoonClient({
    serializationType: SerializationType.JSON,
    wireType: WireType.JSON,
    timeout: 5000,
    url: 'http://localhost:8080/api/v1/events',
    headers: {
        'X-User-ID': 'user-1'
    }
});

//  send the request
raccoonClient
    .send(jsonEvents)
    .then((result) => {
        logger.log('Result:', result);
    })
    .catch((error) => {
        logger.error('Error:', error);
    });
```

Finally run this script using:
```bash
$ node main.js
```

```mdx-code-block
</TabItem>
</Tabs>
```



## Where To Go Next

For more detail about publishing events to Raccoon, you can read the [detailed document](guides/publishing.md) under the guides section. To understand more about how Raccoon works, you can go to the [architecture document](concepts/architecture.md).