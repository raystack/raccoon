---
toc_max_heading_level: 4
---

# Configurations

This page contains reference for all the application configurations for Raccoon.

## Table of Contents

- [Server](configurations.md#server)
- [Worker](configurations.md#worker)
- [Event Distribution](configurations.md#event-distribution)
- [Publisher](configurations.md#publisher)
- [Metric](configurations.md#metric)
- [Log](configurations.md#log)

## Server

### `SERVER_WEBSOCKET_PORT`

Port for the service to listen.

- Type: `Optional`
- Default value: `8080`

### `SERVER_WEBSOCKET_MAX_CONN`

Maximum connection that can be handled by the server instance. You want to set it according to your resource utilization. You also need to check the [limit of open file descriptor allowed](https://docs.oracle.com/cd/E19623-01/820-6168/file-descriptor-requirements.html#:~:text=Linux%20systems%20limit%20the%20number,worker%20threads%20will%20be%20blocked.) by the OS.

- Type: `Optional`
- Default value: `30000`

### `SERVER_WEBSOCKET_READ_BUFFER_SIZE`

Specify I/O buffer sizes in bytes: [Refer gorilla websocket API](https://pkg.go.dev/github.com/gorilla/websocket#hdr-Buffers)

- Type: `Optional`
- Default value: `10240`

### `SERVER_WEBSOCKET_WRITE_BUFFER_SIZE`

Specify I/O buffer sizes in bytes: [Refer gorilla websocket API](https://pkg.go.dev/github.com/gorilla/websocket#hdr-Buffers)

- Type: `Optional`
- Default value: `10240`

### `SERVER_WEBSOCKET_CONN_ID_HEADER`

Unique identifier for the server to maintain the connection. A single uniq id can only connect once in a session. If, there is a subsequence connection with the same uniq id the connection will be rejected.

- Example value: `X-User-ID`
- Type: `Required`

### `SERVER_WEBSOCKET_CONN_GROUP_HEADER`

Additional identifier for the server to maintain the connection. Value of the conn group header combined with user id will act as unique identifier instead of only user id. You can use this if you want to differentiate between user groups or clients e.g(mobile, web). The group names is used as conn_group tag in some of the metrics.

- Example value: `X-User-Group`
- Type: `Optional`

### `SERVER_WEBSOCKET_CONN_GROUP_DEFAULT`

Default connection group name. The default is fallback when `SERVER_WEBSOCKET_CONN_GROUP_HEADER` is not set or when the value of group header is empty. In case the connection group default is clashing with your actual group name, override this config.

- Default value: `--default--`
- Type: `Optional`

### `SERVER_WEBSOCKET_PING_INTERVAL_MS`

Interval of each ping to client. The interval is in seconds.

- Type: `Optional`
- Default value: `30`

### `SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS`

Wait time for client to send Pong message back after server sends Ping. When the time exceeded the connection will be dropped.

- Type `Optional`
- Default value: `60`

### `SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS`

Timeout Deadline set on the writes. On timeout the websocket state is corrupt and all future writes will return error: [Refer gorilla websocket API](https://pkg.go.dev/github.com/gorilla/websocket#Conn.SetWriteDeadline)

- Type `Optional`
- Default value: `5`

### `SERVER_WEBSOCKET_PINGER_SIZE`

Number of goroutine spawn to Ping clients.

- Type `Optional`
- Default value: `1`

### `SERVER_WEBSOCKET_CHECK_ORIGIN`

Toggle CORS check function. Set `true` to check each request origin. Set `false` to disable check origin and allow every request. Check origin function check against `Origin` header.

- Type: `Optional`
- Default value: `true`

### `SERVER_CORS_ENABLED`

The config decides whether to enable the cors middleware and thus allow CORS requests. This config only enables CORS for rest services. For websocket, refer `SERVER_WEBSOCKET_CHECK_ORIGIN`

- Type `Optional`
- Default value: `false`

### `SERVER_CORS_ALLOWED_ORIGIN`

The server decides which origin to allow. The configuration is expected to space separated. Multiple values are supported. The value requires `SERVER_CORS_ENABLED` to be true to take effect. If you want to allow all host headers. You can pass `*` as the value.

- Type `Optional`
- Default Value ``

### `SERVER_CORS_ALLOWED_METHODS`

The http methods allowed when it's a cross origin request. The http methods are expected to be space separated. 

- Type `Optional`
- Default Value `GET HEAD POST OPTIONS`

### `SERVER_CORS_ALLOWED_HEADERS`

The http request headers which are allowed when request is cross origin. The input expects to add any additional headers which is going to be sent by the client ex: `Authorization`. Headers which are essential for the functioning of Raccoon like Content-Type, Connection-Id & Group headers are added by default and need not be passed as configuration.

- Type `Optional`
- Default Value ``

### `SERVER_CORS_ALLOW_CREDENTIALS`

AllowCredentials can be used to specify that the user agent may pass authentication details along with the request. 

- Type `Optional`
- Default Value `false`

### `SERVER_CORS_PREFLIGHT_MAX_AGE_SECONDS`

Replies with a header for clients on how long to cache the response of the preflight request. It's not enforceable. The max value is 600s

- Type `Optional`
- Default Value `0`

### `SERVER_BATCH_DEDUP_IN_CONNECTION_ENABLED`

The server decides whether or not to handle duplicate batches for the active connection. If a batch is sent with a duplicate ReqGUID, the server uses best attempts to discard the duplicate batches. Set `true` to enable the setting.

- Type `Optional`
- Default value: `false`

## Worker

### `WORKER_BUFFER_CHANNEL_SIZE`

Maximum batch that service can handle when workers are busy. When the number of batch is exceeded, the worker will back-pressure causing websocket to stop reading new request.

- Type `Optional`
- Default value: `100`

### `WORKER_BUFFER_FLUSH_TIMEOUT_MS`

Upon shutdown, the worker try to finish processing events in buffer before the timeout exceeded. When the timeout exceeded, the worker is forcefully closed.

- Type `Optional`
- Default value: `5000`

### `WORKER_POOL_SIZE`

No of workers that processes the events concurrently.

- Type `Optional`
- Default value: `5`

### `WORKER_KAFKA_DELIVERY_CHANNEL_SIZE`

Delivery channel is implementation detail where the kafka client asks for channel in the [produce API](https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/producer_example/producer_example.go#L51). The publisher uses the channel to wait for the events to be delivered. The channel contains the status delivery of the events. Normally you won't need to touch this.

- Type `Optional`
- Default value: `10`

## Event Distribution

### `EVENT_DISTRIBUTION_PUBLISHER_PATTERN`

Routes events based on given pattern and [type](https://github.com/raystack/proton/blob/main/raystack/raccoon/Event.proto#L31). The pattern is following [go string format](https://golang.org/pkg/fmt/) with event `type` as second argument. The result of the string format will be the kafka topic target of the event.

For example, you pass `%s-event` as `EVENT_DISTRIBUTION_PUBLISHER_PATTERN`. If you send event with `click` type, your event will be forwarded to `click-event` kafka topic on the configured broker. If you send event with `buy` type, your event will be forwarded to `buy-event`.

You can also route the events to single topic irrespective of the type. To do that you can drop `%s` in the value. For example, provided `mobile-events` as value. All incoming events will be routed to `mobile-events` kafka topic.

- Type `Required`
- Default value: `clickstream-%s-log`

## Publishers

### Common

#### `PUBLISHER_TYPE`

The publisher to use for transmitting events.

Publisher specific configuration follows the pattern `PUBLISHER_${TYPE}_*` where `${TYPE}` is the publisher type in upper case.

- Type `Optional`
- Default value: `kafka`
- Possible values: `kafka`, `pubsub`, `kinesis`

### Kafka

#### `PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS`

Kafka brokers IP address where the events are published.

- Example value: localhost:9092
- Type `Required`

#### `PUBLISHER_KAFKA_CLIENT_ACKS`

Number of replica acknowledgement before it send ack back to service.

- Type `Optional`
- Default value: `-1`

#### `PUBLISHER_KAFKA_CLIENT_RETRIES`

Number of retries in case of failure.

- Type `Optional`
- Default value: `2147483647`

#### `PUBLISHER_KAFKA_CLIENT_RETRY_BACKOFF_MS`

Backoff time on retry.

- Type `Optional`
- Default value: `100`

#### `PUBLISHER_KAFKA_CLIENT_STATISTICS_INTERVAL_MS`

librdkafka statistics emit interval. The application also needs to register a stats callback using rd_kafka_conf_set_stats_cb\(\). The granularity is 1000ms. A value of 0 disables statistics.

- Type `Optional`
- Default value: `0`

#### `PUBLISHER_KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES`

Maximum number of messages allowed on the producer queue. This queue is shared by all topics and partitions.

- Type `Optional`
- Default value: `100000`

#### `PUBLISHER_KAFKA_CLIENT_*`

Kafka client config is dynamically configured. You can see other configurations [here](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md).

The configs are mapped to librdkafka configs by removing the `PUBLISHER_KAFKA_CLIENT_` prefix and replacing underscore with a period.

Internally, this is how it looks
```go title="config/publisher.go"
var dynamicKafkaClientConfigPrefix = "PUBLISHER_KAFKA_CLIENT_"

type publisherKafka struct { /* ... */ }

func (k publisherKafka) ToKafkaConfigMap() *confluent.ConfigMap {
	configMap := &confluent.ConfigMap{}
	for key, value := range viper.AllSettings() {
		if strings.HasPrefix(strings.ToUpper(key), dynamicKafkaClientConfigPrefix) {
			clientConfig := key[len(dynamicKafkaClientConfigPrefix):]
			configMap.SetKey(strings.Join(strings.Split(clientConfig, "_"), "."), value)
		}
	}
	return configMap
}
```

- Type `Optional`
- Default value: see the [reference](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md)

#### `PUBLISHER_KAFKA_FLUSH_INTERVAL_MS`

Upon shutdown, the publisher will try to finish processing events in buffer before the timeout exceeded. When the timeout exceeded, the publisher is forcefully closed.

- Type `Optional`
- Default value: `1000`

### PubSub
#### `PUBLISHER_PUBSUB_CREDENTIALS`

Path to the file containing service account credentials. Defaults to the value of `GOOGLE_APPLICATION_CREDENTIALS` environment variable. This is used to authenticate with Google Cloud Platform.

- Type `Required` (if `PUBLISHER_TYPE=pubsub`, otherwise ignored)

#### `PUBLISHER_PUBSUB_PROJECT_ID`

Destination Google Cloud Project ID. Messages will be transmitted to the PubSub topics under this project.

- Type `Required` (if `PUBLISHER_TYPE=pubsub`, otherwise ignored)

#### `PUBLISHER_PUBSUB_TOPIC_AUTOCREATE`

Whether Raccoon should create a topic if it doesn't exist.

- Type `Optional`
- Default value `false`

#### `PUBLISHER_PUBSUB_TOPIC_RETENTION_MS`

How long PubSub should retain messages in a topic (in milliseconds). Valid values must be between 10 minutes and 31 days.

see [pubsub docs](https://cloud.google.com/pubsub/docs/create-topic) for more information.

- Type `Optional`
- Default value `0`

#### `PUBLISHER_PUBSUB_PUBLISH_DELAY_THRESHOLD_MS`

Maximum time to wait for before publishing a batch of messages.

- Type `Optional`
- Default value `10`

#### `PUBLISHER_PUBSUB_PUBLISH_COUNT_THRESHOLD`

Maximum number of message to accumulate before transmission.

- Type `Optional`
- Default value `100`

#### `PUBLISHER_PUBSUB_PUBLISH_BYTE_THRESHOLD`

Maximum buffer size (in bytes)

- Type `Optional`
- Default value `1000000` (~1MB)

#### `PUBLISHER_PUBSUB_PUBLISH_TIMEOUT_MS`

How long to wait before aborting a publish operation.

- Type `Optional`
- Default value `60000` (1 Minute)

### Kinesis

#### `PUBLISHER_KINESIS_AWS_REGION`

AWS Region of the target kinesis stream. The value of `AWS_REGION` is used as fallback if this variable is not set.

- Type `Required`

#### `PUBLISHER_KINESIS_CREDENTIALS`

Path to [AWS Credentials file](https://docs.aws.amazon.com/sdkref/latest/guide/file-format.html). 

You can also specify the credentials using `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables. 

- Type `Required`

#### `PUBLISHER_KINESIS_STREAM_AUTOCREATE`

Whether Raccoon should create a stream if it doesn't exist.

NOTE: We recommend that you create all streams that you need to publish to ahead of time.

- Type `Optional`
- Default value `false`

#### `PUBLISHER_KINESIS_STREAM_MODE`

This configuration variable controls the `StreamMode` of the
streams created by Raccoon.

- Type `Optional`
- Default value `ON_DEMAND`
- Possible values: `ON_DEMAND`, `PROVISIONED`

#### `PUBLISHER_KINESIS_STREAM_SHARDS`

This controls the number of shards configured for a stream created by Raccoon.

- Type `Optional`
- Default value `4`

#### `PUBLISHER_KINESIS_STREAM_PROBE_INTERVAL_MS`

This specifies the time delay between stream status checks.

- Type `Optional`
- Default value `1000`


#### `PUBLISHER_KINESIS_PUBLISH_TIMEOUT_MS`

How long to wait for before aborting a publish operation.

- Type `Optional`
- Default value `60000`

## Metric

### `METRIC_RUNTIME_STATS_RECORD_INTERVAL_MS`

The time interval between recording runtime stats of the application in the instrumentation. It's recommended to keep this value equivalent to flush interval when using statsd and your collector's scrape interval when using prometheus as your instrumentation.

- Type `Optional`
- Default Value: `10000`

### `METRIC_STATSD_ENABLED`

Flag to enable export of statsd metric

- Type `Optional`
- Default value: `false`

### `METRIC_STATSD_ADDRESS`

Address to reports the service metrics.

- Type `Optional`
- Default value: `:8125`

### `METRIC_STATSD_FLUSH_PERIOD_MS`

Interval for the service to push metrics.

- Type `Optional`
- Default value: `10000`

### `METRIC_PROMETHEUS_ENABLED`

Flag to enable a prometheus http server to expose metrics.

- Type `Optional`
- Default value: `false`

### `METRIC_PROMETHEUS_PATH`

The path at which prometheus server should serve metrics.

- Type `Optional`
- Default value: `/metrics`

### `METRIC_PROMETHEUS_PORT`

The port number on which prometheus server will be listening for metric scraping requests.

- Type `Optional`
- Default value: `9090`

## Log

### `LOG_LEVEL`

Level available are `info` `panic` `fatal` `error` `warn` `info` `debug` `trace`.

- Type `Optional`
- Default value: `info`

## Event

### `EVENT_ACK`

Based on this parameter the server decides when to send the acknowledgement to the client. Supported values are `0` and `1`.

- Type `Optional`
- Default value: `0`
