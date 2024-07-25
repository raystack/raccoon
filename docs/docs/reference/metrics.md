---
toc_max_heading_level: 4
---
# Metrics

Raccoon supports `statsd` and `prometheus` as two ways to report metrics. For statsd, we recommend using [Telegraf](https://github.com/influxdata/telegraf) as a collection agent.

This page contains the reference for all the metrics exposed by Raccoon.

## Table of Contents

- [Server Connection](metrics.md#server-connection)
- [Publisher](metrics.md#publisher)
  - [Kafka](#kafka)
  - [PubSub](#pubsub)
  - [Kinesis](#kinesis)
- [Event Delivery](metrics.md#event-delivery)

## Server Connection

### `server_ping_failure_total`

Total ping that server fails to send

- Type: `Counting`
- Tags: `conn_group=*`

### `server_pong_failure_total`

Total pong that server fails to send

- Type: `Counting`
- Tags: `conn_group=*`

### `connections_count_current`

Number of alive connections

- Type: `Gauge`
- Tags: `conn_group=*`

### `user_connection_success_total`

Number of successful connections established to the server

- Type: `Count`
- Tags: `conn_group=*`

### `user_connection_failure_total`

Number of fail connections established to the server

- Type: `Count`
- Tags: `reason=ugfailure` `reason=exists` `reason=serverlimit` `conn_group=*`

### `user_session_duration_milliseconds`

Duration of alive connection per session per connection

- Type: `Timing`
- Tags: `conn_group=*`

### `conn_close_err_count`

Number of connection close errors encountered

- Type: `Count`
- Tags: NA

## Publisher
### Kafka
#### `kafka_messages_delivered_total`

Number of events delivered to Kafka. 

- Type: `Count`
- Tags: `topic=topicname` `conn_group=*` `event_type=*`

#### `kafka_messages_undelivered_total`

Number of events not delivered to Kafka.

- Type: `Count`
- Tags: `topic=topicname` `conn_group=*` `event_type=*`


#### `kafka_unknown_topic_failure_total`

Number of delivery failure caused by topic does not exist in kafka.

- Type: `Count`
- Tags: `topic=topicname` `event_type=*`

#### `kafka_tx_messages_total`

Total number of messages transmitted \(produced\) to Kafka brokers.

- Type: `Gauge`

#### `kafka_tx_messages_bytes_total`

Total number of message bytes \(including framing, such as per-Message framing and MessageSet/batch framing\) transmitted to Kafka brokers

- Type: `Gauge`

#### `kafka_brokers_tx_total`

Total number of requests sent to Kafka brokers

- Type: `Gauge`
- Tags: `broker=broker_nodes`

#### `kafka_brokers_tx_bytes_total`

Total number of bytes transmitted to Kafka brokers

- Type: `Gauge`
- Tags: `broker=broker_nodes`

#### `kafka_brokers_rtt_average_milliseconds`

Broker latency / round-trip time in microseconds

- Type: `Gauge`
- Tags: `broker=broker_nodes`

#### `ack_event_rtt_ms`

Time taken from ack function called by kafka producer to processed by the ack handler.

- Type: `Timing`
- Tags: NA

#### `event_rtt_ms`

Time taken from event is consumed from the queue to be acked by the ack handler.

- Type: `Timing`
- Tags: NA

#### `kafka_producebulk_tt_ms`

Response time of produce batch method of the kafka producer

- Type `Timing`
- Tags: NA


### PubSub

#### `pubsub_messages_delivered_total`

Number of events delivered to PubSub. 

- Type: `Count`
- Tags: `topic=topicname` `conn_group=*` `event_type=*`

#### `pubsub_messages_undelivered_total`

Number of events that were not delivered to PubSub.

- Type: `Count`
- Tags: `topic=topicname` `conn_group=*` `event_type=*`


#### `pubsub_unknown_topic_failure_total`

Number of delivery failures caused by non-existence of topic in PubSub.

- Type: `Count`
- Tags: `topic=topicname` `event_type=*` `conn_group=*`

#### `pubsub_topic_throughput_exceeded_total`

Number of delivery failures caused by exceeding throughput limits on PubSub.

- Type: `Count`
- Tags: `topic=topicname` `event_type=*` `conn_group=*`

#### `pubsub_topics_limit_exceeded_total`

Number of delivery failures caused by exceeding the limit on number of Topics on PubSub.

- Type: `Count`
- Tags: `topic=topicname` `event_type=*` `conn_group=*`

### Kinesis

#### `kinesis_messages_delivered_total`

Number of events successfully delivered to Kinesis. 

- Type: `Count`
- Tags: `stream=streamname` `conn_group=*` `event_type=*`

#### `kinesis_messages_undelivered_total`

Number of events not delivered to Kinesis.

- Type: `Count`
- Tags: `stream=streamname` `conn_group=*` `event_type=*`


#### `kinesis_unknown_stream_failure_total`

Number of delivery failures caused by non-existence of stream in Kinesis.

- Type: `Count`
- Tags: `stream=streamname` `event_type=*` `conn_group=*`

#### `kinesis_stream_throughput_exceeded_total`

Number of delivery failures caused by exceeding shard throughput limits. This error can also occur if the message size of an event exceeds message size limit (1MiB as of the day of this writing). See [Limits and Quotas on Kinesis](https://docs.aws.amazon.com/streams/latest/dev/service-sizes-and-limits.html)

- Type: `Count`
- Tags: `stream=streamname` `event_type=*` `conn_group=*`

#### `kinesis_streams_limit_exceeded_total`

Number of delivery failures caused due to too many streams in `CREATING` mode. AWS Kinesis limits how many stream creation requests can be submitted in parallel to 5.

- Type: `Count`
- Tags: `stream=streamname` `event_type=*` `conn_group=*`

## Resource Usage

### `server_mem_gc_triggered_current`

The time the last garbage collection finished in Unix timestamp

- Type: `Gauge`

### `server_mem_gc_pauseNs_current`

Circular buffer of recent GC stop-the-world in Unix timestamp

- Type: `Gauge`

### `server_mem_gc_count_current`

The number of completed GC cycle

- Type: `Gauge`

### `server_mem_gc_pauseTotalNs_current`

The cumulative nanoseconds in GC stop-the-world pauses since the program started

- Type: `Gauge`

### `server_mem_heap_alloc_bytes_current`

Bytes of allocated heap objects

- Type: `Gauge`

### `server_mem_heap_inuse_bytes_current`

HeapInuse is bytes in in-use spans

- Type: `Gauge`

### `server_mem_heap_objects_total_current`

Number of allocated heap objects

- Type: `Gauge`

### `server_go_routines_count_current`

Number of goroutine spawn in a single flush

- Type: `Gauge`

### `server_mem_stack_inuse_bytes_current`

Bytes in stack spans

- Type: `Gauge`

## Event Delivery

Following metrics are event delivery reports. Each metrics reported at a different point in time. See the diagram below for to understand when each metrics are reported.

![Diagram](/assets/metrics_report_time.svg)

### `events_rx_bytes_total`

Total byte received in requests

- Type: `Count`
- Tags: `conn_group=*` `event_type=*`

### `events_rx_total`

Number of events received in requests

- Type: `Count`
- Tags: `conn_group=*` `event_type=*`

### `events_duplicate_total`

Number of duplicate events

- Type: `Count`
- Tags: `conn_group=*` `reason=*`

### `batches_read_total`

Request count

- Type: `Count`
- Tags: `status=failed` `status=success` `reason=*` `conn_group=*`

### `batch_idle_in_channel_milliseconds`

Duration from when the request is received to when the request is processed. High value of this metric indicates the publisher is slow.

- Type: `Timing`
- Tags: `worker=worker-name`

### `event_processing_duration_milliseconds`

Duration from the time request is sent to the time events are published. This metric is calculated per event by following formula `(PublishedTime - SentTime)/CountEvents`

- Type: `Timing`
- Tags: `conn_group=*`

### `server_processing_latency_milliseconds`

Duration from the time request is received to the time events are published. This metric is calculated per event by following formula`(PublishedTime - ReceivedTime)/CountEvents`

- Type: `Timing`
- Tags: `conn_group=*`

### `worker_processing_duration_milliseconds`

Duration from the time request is processed to the time events are published. This metric is calculated per event by following formula`(PublishedTime - ProcessedTime)/CountEvents`

- Type: `Timing`
