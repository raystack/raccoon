# Introduction

Raccoon was built with a primary purpose to source (or collect) user behaviour data in near-real time. User behaviour data is a stream of events that occur when users traverse through a mobile app or website. Raccoon powers analytics systems, big data pipelines and other disparate consumers by providing high volume, high throughput ingestion APIs consuming real time data. Raccoon’s key architecture principle is a realization of an event agnostic backend (accepts events of any type without the type awareness). 
It is this capability that enables Raccoon to evolve into a strong player in the ingestion/collector ecosystem that has real time streaming/analytical needs.

## Architecture

Raccoon written in [GO](https://github.com/golang) is a high throughput, low-latency service that provides an API to ingest streaming data from mobile apps, sites and publish it to Kafka. Raccoon uses the Websocket protocol for peer-to-peer communication, providing long persistent connections, with no overhead of additional headers sizes as in http protocol. Protobuf is used as the serialization format that reduces the payload sizes further. It provides an event type agnostic API that accepts a batch (array) of events in protobuf format. Refer [here](https://github.com/odpf/proton/tree/main/odpf/raccoon) for proto definition format that Raccoon accepts.


## System Design

<p align="center"><img src="../assets/raccoon_hld.png" /></p>


At a high level, the following sequence details the architecture.
* Clients make websocket connections to Raccoon by performing a http GET API call, with headers to upgrade to websocket.
* Raccoon uses [gorilla websocket](https://github.com/gorilla/websocket) handlers and for each connection the handlers spawn a goroutine to handle incoming requests.
* After the connection has been established, clients can send the events. 
* Clients can send the request anytime as long as the websocket connection is alive.
* When the events are available to be consumed, the handler will deserialize the proto events and forward them to the buffered channel.
* A pool of worker go routines works off the buffered channel
* Each worker iterates over the events' batch, determines the topic based on the type and serializes the bytes to the Kafka producer synchronously.

Note: The internals of each of the components like channel size, buffer sizes, publisher properties etc., are configurable enabling Raccoon to be provisioned according to the system/event characteristics and load.

### Connections
Raccoon has long running persistent connections with the client. Once a client makes a http request with a websocket upgrade header, raccoon upgrades the http request to a websocket connection end of which a persistent connection is established with the client.

The following sequence outlines the connection handling by Raccoon.
* Fetch connection id details from the initial request header based on the configured header name in `SERVER_WEBSOCKET_CONN_UNIQ_ID_HEADER`. The header name uniquely identifies a client. A client in this case can be the user in the app. There can be multiple connections from the same client. The no., of connections allowed per client is determined by `SERVER_WEBSOCKET_MAX_CONN`.
* Once the connection id is fetched, verify if the user has connection limit reached based on the configured `SERVER_WEBSOCKET_MAX_CONN`. For each client an internal map stores the `SERVER_WEBSOCKET_MAX_CONN` along with the connection objects. On reaching the max connections for the client, the connection is disconnected with an appropriate error message as a response proto.
* Upgrade the connection
* Add this user-id -> connection mapping
* Add ping/pong handlers on this connection, readtimeout deadline. More about these handlers in the following sections
* Handle the message and send it to the events-channel
* Remove connection/user when the client closes the connection

### Acknowledging events

Event acknowledgements was designed to signify if the events batch is received and sent to Kafka successfully. This will enable the clients to retry on failed event delivery.
However Raccoon chooses to send the acknowledgments as soon as it receives and deserializes the events successfully using the proto `EventRequest`. The acks are sent even before it is produced to Kafka. The following picture depicts the sequence of the event ack.

<p align="center"><img src="../assets/raccoon_ack.png" /></p>

Pros:

- Performant as it does not wait for kafka/network round trip for each batch of events.

Cons:

- Potential data-loss and the clients do not get a chance to retry/resend the events. The possiblity of data-loss occurs when the kafka borker cluster is facing a downtime. 

Considering that kafka is set up in a clustered, cross-region, cross-zone environment, the chances of it going down are mostly unlikely. In case if it does, the amount of events lost is negligible considering it is a streaming system and is expected to forward millions of events/sec.

When an [EventRequest](https://github.com/odpf/proton/blob/main/odpf/raccoon/EventRequest.proto) proto below containing events are sent over the wire 

```
message EventRequest {
  //unique guid generated by the client for this request
  string req_guid = 1;
  // time probably when the client sent it
  google.protobuf.Timestamp sent_time = 2;
  // actual events 
  repeated Event events = 3;
}
```

a corresponding [EventResponse](https://github.com/odpf/proton/blob/main/odpf/raccoon/EventResponse.proto) is sent by the server on the same connection that the events were consumed.

```
message EventResponse {
  Status status = 1;
  Code code = 2;
      /* time when the response is generated */
  int64 sent_time = 3;
      /* failure reasons if any */
  string reason = 4;
      /* Usually detailing the success/failures */
  map<string, string> data = 5;
}
```

### Event Distribution

Event distribution works by finding the type for each event in the batch and sending them to appropriate kafka topic. The topic name is determined by the following code 

```
topic := fmt.Sprintf(pr.topicFormat, event.Type)
```

where 
**topicformat** - is the configured pattern `EVENT_DISTRIBUTION_PUBLISHER_PATTERN`
**type** - is the type set by the client when the event proto is generated

For eg. setting the 
```
EVENT_DISTRIBUTION_PUBLISHER_PATTERN=topic-%s-log
```
and a type such as ```type=viewed``` in the [event](https://github.com/odpf/proton/blob/main/odpf/raccoon/Event.proto) format

```
message Event {
  /*
  `eventBytes` is where you put bytes serialized event.
  */
  bytes eventBytes = 1;
  /*
  `type` denotes an event type that the producer of this proto message may set.
  It is currently used by raccoon to distribute events to respective Kafka topics. However the
  users of this proto can use this type to set strings which can be processed in their
  ingestion systems to distribute or perform other functions.
  */
  string type = 2;
 }
```

will have the event sent to a topic like
```topic-viewed-log```

The event distribution does not depend on any partition logic. So events can be randomnly distrbuted to any kafka partition.

### Event Deserialization

The top level wrapper `EventRequest` is deserialized which provides a list of events of type `Event` proto. This event wrapper composes of serialized bytes, which is the actual event, set in the field `bytes` inside the `Event` proto. Raccoon does not open this underlying bytes. The deserialization is used to unwrap the event type and determine the topic that the `eventBytes` (an event) need to be sent to.

### Channels
Buffered Channels are used to store the incoming events' batch. The server acknowledges the client's receiving message. The channel sizes can be configured based on the load & capacity.

### Keeping connections alive
The server ensures that the connections are recyclable. It adopts mechanisms to check connection time idleness. The handlers ping clients very 30 seconds (configurable). If the client does not respond within a stipulated time the connection is marked as corrupt. Every subsequent read/write message there after on this connection fails. Raccoon removes the connections post this.
Clients can also ping the server while the server responds with pongs to these pings. Clients can programmtically reconnect on failed or corrupt server connections.

## Components

### Kafka producer
Raccoon uses [confluent go kafka](https://github.com/confluentinc/confluent-kafka-go) as the producer client to publish events. Publishing events are light weight and relies on kafka producer's retries. Confluent internally uses librdkafka which produces events asynchronously. Application writes messages using a functional based producer API

`Produce(message, deliveryChannel)`
 -- `deliveryChannel` is where the delivery reports or acknowledgements are received.

Raccoon internally checks for these delivery reports before pulling the next batch of events. On failed deliveries the appropriate metrics are updated. This mechanism makes the events delivery synchronous and a reliable events delivery.

### Observability Stack
Raccoon internally uses [statsd](https://gopkg.in/alexcesaro/statsd.v2) go module client to export metrics in StatsD line protocol format. A recommended choice for observability stack would be to host [telegraf](https://www.influxdata.com/time-series-platform/telegraf/) as the receiver of these measurements and expoert it to [influx](https://www.influxdata.com/get-influxdb/), influx to store the metrics, [grafana](https://grafana.com/) to build dashboards using Influx as the source.