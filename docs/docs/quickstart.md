import Tabs from '@theme/Tabs'
import TabItem from '@theme/TabItem'
import CodeBlock from '@theme/CodeBlock'

# Quickstart

This document will guide you on how to get Raccoon + Kafka setup running locally. This document assumes that you have Docker (with Docker Compose) and Kafka installed on your system. 

## Run Raccoon with Docker Compose

Here's a minimal setup that runs a single node kafka-cluster along with raccoon:

```yaml title="docker-compose.yml"
networks:
  raccoon-network:

services:
  zookeeper:
    image: confluentinc/cp-zookeeper:5.1.2
    hostname: zookeeper
    container_name: zookeeper
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - raccoon-network

  kafka:
    image: confluentinc/cp-kafka:5.1.2
    hostname: kafka
    container_name: kafka
    depends_on:
      - zookeeper
    ports:
      - "9094:9094"
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_ADVERTISED_LISTENERS: INSIDE://kafka:9092,OUTSIDE://localhost:9094
      KAFKA_LISTENERS: INSIDE://:9092,OUTSIDE://:9094
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: INSIDE
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_SOCKET_REQUEST_MAX_BYTES: 1000000000
      CONFLUENT_METRICS_ENABLE: 'false'
    links:
      - zookeeper
    networks:
      - raccoon-network
  raccoon:
    image: raystack/raccoon
    hostname: raccoon
    container_name: raccoon
    stdin_open: true
    tty: true
    depends_on:
      - kafka
    environment:
      SERVER_WEBSOCKET_PORT: "8080"
      SERVER_WEBSOCKET_CHECK_ORIGIN: "true"
      SERVER_CORS_ENABLED: "true"
      SERVER_CORS_ALLOWED_ORIGIN: "http://localhost:3000 http://localhost:8080"
      SERVER_CORS_ALLOWED_METHODS: "GET HEAD POST OPTIONS"
      SERVER_WEBSOCKET_CONN_ID_HEADER: "X-User-ID"
      SERVER_WEBSOCKET_CONN_GROUP_HEADER: "X-User-Group"
      SERVER_GRPC_PORT: 8081
      EVENT_DISTRIBUTION_PUBLISHER_PATTERN: "event-log"
      PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS: "kafka:9092"
    ports:
      - "8080:8080"
      - "8081:8081"
    networks:
      - raccoon-network
```

This setup is configured to publish all events to `event-log` topic. You can also configure Raccoon to [route events to different topics based on the event type.](concepts/architecture.md#event-distribution)

Copy the file to your local system and run the following to start Raccoon.
```bash
$ docker compose up
```

To test whether Raccoon is running or not, you can try to ping the server.

```bash
$ curl http://localhost:8080/ping
```

To verify the event published by Raccoon. First, you need to start a Kafka listener. In a seperate terminal run:
```bash
$ kafka-console-consumer --bootstrap-server localhost:9094 --topic 'event-log'
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
<TabItem value='java'>
```

Make sure you have Java installed. Raccoon java client requires:
* `JDK` version >= 8
* `Gradle` version >= 7

Begin by creating a new java project in a folder called `java-raccoon-example`
```bash
$ mkdir java-raccoon-example
$ cd java-raccoon-example
$ gradle init --type=java-application
```

Add `io.odpf.raccoon` version `0.1.5-rc` in your `build.gradle`. It should look something like this:
```groovy
plugins {
    // Apply the application plugin to add support for building a CLI application in Java.
    id 'application'
}

repositories {
    // Use Maven Central for resolving dependencies.
    mavenCentral()
}

dependencies {
    // Use JUnit test framework.
    testImplementation libs.junit

    // This dependency is used by the application.
    implementation libs.guava
    
    // Raccoon Client library
    implementation group: 'io.odpf', name: 'raccoon', version: '0.1.5-rc'
}

// Apply a specific Java toolchain to ease working on different environments.
java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

application {
    // Define the main class for the application.
    mainClass = 'org.example.App'
}
```

Edit the `App.java` file and add the following code:
```java title=src/java/org/example/App.java showLineNumbers
package org.example;

import io.odpf.raccoon.client.RestConfig;
import io.odpf.raccoon.client.RaccoonClient;
import io.odpf.raccoon.client.RaccoonClientFactory;
import io.odpf.raccoon.model.Event;
import io.odpf.raccoon.model.Response;
import io.odpf.raccoon.model.ResponseStatus;
import io.odpf.raccoon.serializer.JsonSerializer;
import io.odpf.raccoon.wire.ProtoWire;

public class App {

    public static void main(String[] args) {
        RestConfig config = RestConfig.builder()
                  .url("http://localhost:8080/api/v1/events")
                  .header("x-user-id", "123")
                  .serializer(new JsonSerializer()) // default is Json
                  .marshaler(new ProtoWire()) // default is Json
                  .retryMax(5) // default is 3
                  .retryWait(2000) // default is one second
                  .build();

        // get the rest client instance.
        RaccoonClient Client = RaccoonClientFactory.getRestClient(config);

        Response res = Client.send(new Event[]{
                new Event("page", "EVENT".getBytes())
        });

        if (res.isSuccess() && res.getStatus() == ResponseStatus.STATUS_SUCCESS) {
                System.out.println("The event was sent successfully");
        }
    }
}

```

Run the application using
```bash
$ gradle run
```

```mdx-code-block
</TabItem>
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

Create a `main.mjs` file with the following contents:
```js title="main.mjs" showLineNumbers
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
$ node main.mjs
```

```mdx-code-block
</TabItem>
</Tabs>
```


## Where To Go Next

For more detail about publishing events to Raccoon, you can read the [detailed document](guides/publishing.md) under the guides section. To understand more about how Raccoon works, you can go to the [architecture document](concepts/architecture.md).