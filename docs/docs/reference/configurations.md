# Configurations

This page contains reference for all the application configurations for Raccoon.

## Table of Contents

* [Server](configurations.md#server)
* [Worker](configurations.md#worker)
* [Event Distribution](configurations.md#event-distribution)
* [Publisher](configurations.md#publisher)
* [Metric](configurations.md#metric)
* [Log](configurations.md#log)

## Server

### `SERVER_WEBSOCKET_PORT`

Port for the service to listen.

* Type: `Optional`
* Default value: `8080`

### `SERVER_WEBSOCKET_MAX_CONN`

Maximum connection that can be handled by the server instance. You want to set it according to your resource utilization. You also need to check the [limit of open file descriptor allowed](https://docs.oracle.com/cd/E19623-01/820-6168/file-descriptor-requirements.html#:~:text=Linux%20systems%20limit%20the%20number,worker%20threads%20will%20be%20blocked.) by the OS.

* Type: `Optional`
* Default value: `30000`

### `SERVER_WEBSOCKET_READ_BUFFER_SIZE`

Specify I/O buffer sizes in bytes: [Refer gorilla websocket API](https://pkg.go.dev/github.com/gorilla/websocket#hdr-Buffers)

* Type: `Optional`
* Default value: `10240`

### `SERVER_WEBSOCKET_WRITE_BUFFER_SIZE`

Specify I/O buffer sizes in bytes: [Refer gorilla websocket API](https://pkg.go.dev/github.com/gorilla/websocket#hdr-Buffers)

* Type: `Optional`
* Default value: `10240`

### `SERVER_WEBSOCKET_CONN_ID_HEADER`

Unique identifier for the server to maintain the connection. A single uniq id can only connect once in a session. If, there is a subsequence connection with the same uniq id the connection will be rejected.

* Example value: `X-User-ID`
* Type: `Required`

### `SERVER_WEBSOCKET_CONN_GROUP_HEADER`

Additional identifier for the server to maintain the connection. Value of the conn group header combined with user id will act as unique identifier instead of only user id. You can use this if you want to differentiate between user groups or clients e.g(mobile, web). The group names is used as conn_group tag in some of the metrics.

* Example value: `X-User-Group`
* Type: `Optional`

### `SERVER_WEBSOCKET_CONN_GROUP_DEFAULT`

Default connection group name. The default is fallback when `SERVER_WEBSOCKET_CONN_GROUP_HEADER` is not set or when the value of group header is empty. In case the connection group default is clashing with your actual group name, override this config.

* Default value: `--default--`
* Type: `Optional`

### `SERVER_WEBSOCKET_PING_INTERVAL_MS`

Interval of each ping to client. The interval is in seconds.

* Type: `Optional`
* Default value: `30`

### `SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS`

Wait time for client to send Pong message back after server sends Ping. When the time exceeded the connection will be dropped.

* Type `Optional`
* Default value: `60`

### `SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS`

Timeout Deadline set on the writes. On timeout the websocket state is corrupt and all future writes will return error: [Refer gorilla websocket API](https://pkg.go.dev/github.com/gorilla/websocket#Conn.SetWriteDeadline)

* Type `Optional`
* Default value: `5`

### `SERVER_WEBSOCKET_PINGER_SIZE`

Number of goroutine spawn to Ping clients.

* Type `Optional`
* Default value: `1`

### `SERVER_WEBSOCKET_CHECK_ORIGIN`

Toggle CORS check function. Set `true` to check each request origin. Set `false` to disable check origin and allow every request. Check origin function check against `Origin` header.

* Type: `Optional`
* Default value: `true`

### `SERVER_BATCH_DEDUP_IN_CONNECTION_ENABLED`

The server decides whether or not to handle duplicate batches for the active connection. If a batch is sent with a duplicate ReqGUID, the server uses best attempts to discard the duplicate batches. Set `true` to enable the setting.

* Type `Optional`
* Default value: `false`

## Worker

### `WORKER_BUFFER_CHANNEL_SIZE`

Maximum batch that service can handle when workers are busy. When the number of batch is exceeded, the worker will backpressure causing websocket to stop reading new request.

* Type `Optional`
* Default value: `100`

### `WORKER_BUFFER_FLUSH_TIMEOUT_MS`

Upon shutdown, the worker try to finish processing events in buffer before the timeout exceeded. When the timeout exceeded, the worker is forcefully closed.

* Type `Optional`
* Default value: `5000`

### `WORKER_POOL_SIZE`

No of workers that processes the events concurrently.

* Type `Optional`
* Default value: `5`

### `WORKER_KAFKA_DELIVERY_CHANNEL_SIZE`

Delivery channel is implementation detail where the kafka client asks for channel in the [produce API](https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/producer_example/producer_example.go#L51). The publisher uses the channel to wait for the events to be delivered. The channel contains the status delivery of the events. Normally you won't need to touch this.

* Type `Optional`
* Default value: `10`

## Event Distribution

### `EVENT_DISTRIBUTION_PUBLISHER_PATTERN`

Routes events based on given pattern and [type](https://github.com/goto/proton/blob/main/goto/raccoon/Event.proto#L31). The pattern is following [go string format](https://golang.org/pkg/fmt/) with event `type` as second argument. The result of the string format will be the kafka topic target of the event.

For example, you pass `%s-event` as `EVENT_DISTRIBUTION_PUBLISHER_PATTERN`. If you send event with `click` type, your event will be forwareded to `click-event` kafka topic on the configured broker. If you send event with `buy` type, your event will be forwarded to `buy-event`.

You can also route the events to single topic irrespective of the type. To do that you can drop `%s` in the value. For example, provided `mobile-events` as value. All incoming events will be routed to `mobile-events` kafka topic.

* Type `Required`
* Default value: `clickstream-%s-log`

## Publisher

### `PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS`

Kafka brokers IP address where the events are published.

* Example value: localhost:9092
* Type `Required`

### `PUBLISHER_KAFKA_CLIENT_ACKS`

Number of replica acknowledgement before it send ack back to service.

* Type `Optional`
* Default value: `-1`

### `PUBLISHER_KAFKA_CLIENT_RETRIES`

Number of retries in case of failure.

* Type `Optional`
* Default value: `2147483647`

### `PUBLISHER_KAFKA_CLIENT_RETRY_BACKOFF_MS`

Backoff time on retry.

* Type `Optional`
* Default value: `100`

### `PUBLISHER_KAFKA_CLIENT_STATISTICS_INTERVAL_MS`

librdkafka statistics emit interval. The application also needs to register a stats callback using rd\_kafka\_conf\_set\_stats\_cb\(\). The granularity is 1000ms. A value of 0 disables statistics.

* Type `Optional`
* Default value: `0`

### `PUBLISHER_KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES`

Maximum number of messages allowed on the producer queue. This queue is shared by all topics and partitions.

* Type `Optional`
* Default value: `100000`

### `PUBLISHER_KAFKA_CLIENT_*`

Kafka client config is dynamically configured. You can see for other configuration [here](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md)

* Type `Optional`
* Default value: see the [reference](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md)

### `PUBLISHER_KAFKA_FLUSH_INTERVAL_MS`

Upon shutdown, the publisher will try to finish processing events in buffer before the timeout exceeded. When the timeout exceeded, the publisher is forcefully closed.

* Type `Optional`
* Default value: `1000`

## Metric

### `METRIC_STATSD_ADDRESS`

Address to reports the service metrics.

* Type `Optional`
* Default value: `:8125`

### `METRIC_STATSD_FLUSH_PERIOD_MS`

Interval for the service to push metrics.

* Type `Optional`
* Default value: `10000`

## Log

### `LOG_LEVEL`

Level available are `info` `panic` `fatal` `error` `warn` `info` `debug` `trace`.

* Type `Optional`
* Default value: `info`

## Event

### `EVENT_ACK`

Based on this parameter the server decides when to send the acknowledgement to the client. Supported values are `0` and `1`.

* Type `Optional`
* Default value: `0`