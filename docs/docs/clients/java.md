# Java

## Requirements
Make sure you have Java JDK `>=8` and Gradle `>=7` installed on your system. See installation instructions for [openjdk](https://openjdk.org/install/) and [gradle](https://gradle.org/install/) for more information.
## Installation

In your `build.gradle` file add `io.odpf.raccoon` as a dependency.
```groovy
dependencies {
    implementation group: 'io.odpf', name: 'raccoon', version: '0.1.5-rc'
}
```

## Usage

### Quickstart

Below is a self contained example of Raccoon's Java client that uses the REST API to publish events
```java title="App.java"
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

### Guide

#### Creating a client

Raccoon's Java client only supports sending events over Raccoon's HTTP/JSON (REST) API.

To create a client, you must pass the `io.odpf.raccoon.client.RestConfig` object to the client factory `io.odpf.raccoon.client.RaccoonClientFactory.getRestClient()`.

You can use `RestConfig.builder()` as a convenient way of building the config object.

Here's a minimal example of what it looks like:
```java
RestConfig config = RestConfig.builder()
        .url("http://localhost:8080/api/v1/events")
        .build();
RacconClient client = RaccoonClientFactory.getRestClient(config);
```

#### Sending events

Events can be sent to raccoon use `RestClient.send()` method. The `send()` methods accepts an array of `io.odpf.raccoon.model.Event`

Here's a minimal example of what this could look like:
```java
Event[] events = new Event[]{
    new Event("event_type", obj)
};
client.send(events);
```

Each event has a `type` and `data` field. `type` denotes the event type. This is used by raccoon to route the event to a specific topic downstream. `data` field contains the payload. This data is serialised by the `serializer` that's configured on the client. 

The following table lists which serializer to use for a given payload type.

| Message Type | Serializer |
| --- | --- |
| JSON | `io.odpf.raccoon.serializer.JsonSerializer` |
| Protobuf | `io.odpf.raccoon.serializer.ProtoSerializer`|

Once a client is constructed with a specific kind of serializer, you may only pass it events of that specific type. In particular, for `JSON` serialiser the event data must be a Java object. While for `PROTOBUF` serialiser the event data must be a protobuf message object
