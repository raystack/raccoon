---
toc_max_heading_level: 4
---

# Troubleshooting

## Scale Up Racoon

Internally, Raccoon has 3 main components that affect the capacity. The server, worker, and publisher. Each component has configurations that can be tune if necessary. Since those 3 components are forming a pipe, you need to make sure none of the components become a bottleneck.

To know the right configuration, you need to simulate with similar throughput as production. You can tune the configuration accordingly.

Following are details of what you can tune.

### Server

Raccoon is using WebSocket as a communication protocol from client to server. Websocket requires maintaining long-running connections. Each connection costs the OS an open file descriptor. When you reach the limit of the configured open file descriptor, the server won't be able to accept a new connection. By default, OS limit the number of the open file descriptor. You can look up how to increase the max open file descriptor. On Unix, you can do `ulimit -n` to check max open file descriptor and `ulimit -n <number>` to set a new limit.

Apart from OS configuration, there are configurations you can tune on Raccoon:

- [SERVER_WEBSOCKET_MAX_CONN](reference/configurations.md#server_websocket_max_conn) To limit Raccoon resource utilization, we enforce a limit on WebSocket connection. The default value is 30000; adjust it if necessary.

### Worker

After the request is deserialized, the server puts the events on the buffer channel. The worker process events from the channel and publishes them downstream. You can think of the worker and the channel as a buffer in case the publisher slows down temporarily.

- [WORKER_BUFFER_CHANNEL_SIZE](reference/configurations.md#worker_buffer_channel_size) Buffer before the events get processed. The more the size, the longer it can tolerate a temporary spike or slow down.
- [WORKER_POOL_SIZE](reference/configurations.md#worker_pool_size) The worker will call the publisher client and wait synchronously. Increase this according to the throughput.

### Publisher

Raccoon has support for `kafka`, `pubsub` and `kinesis` publishers.

#### Kafka
Currently, Raccoon is using [Librd Kafka client Go wrapper](https://github.com/confluentinc/confluent-kafka-go) as publisher client. There is plenty of guides out there to tune Kafka producer. Here are some configurations you can tune.

- [PUBLISHER_KAFKA_CLIENT_BATCH_NUM_MESSAGES](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md)
- [PUBLISHER_KAFKA_CLIENT_ACKS](reference/configurations.md#publisher_kafka_client_acks)
- [PUBLISHER\_KAFKA\_CLIENT_${conf}](reference/configurations.md#publisher_kafka_client_) You can put any [librd kafka configuration](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md) by replacing `${conf}` with upper case'd configuration key and changing the delimiter to underscore. For example, to use `log.queue=true`, you can set `PUBLISHER_KAFKA_CLIENT_LOG_QUEUE=true`

#### PubSub
Raccoon uses [cloud.google.com/go/pubsub](https://pkg.go.dev/cloud.google.com/go/pubsub) as the producer client for publishing events to Google Cloud PubSub.

The default quota limits for writes are:
* 4GiB/second for large regions
* 800MiB/second for medium regions
* 200MiB/second for small regions

A single message (event) must not be bigger 10MiB. Although this limit can be increased by submitting a quota increase request.

Since PubSub is a managed service, you generally only need to worry about hitting quotas or rate limits. Refer to [PubSub documentation](https://cloud.google.com/pubsub/quotas) for more information. 


#### Kinesis

Raccoon uses [github.com/aws/aws-sdk-go-v2/service/kinesis](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/kinesis) as the producer client for publishing events to AWS Kinesis.

AWS Kinesis Data Stream come into two modes:
* Provisioned
* On-Demand

With Provisioned mode, your throughput is computed using the formlua:
```
Throughput/Second = Number of Shards * 1MiB
Records/Second = Number of Shards * 1000
```
Shards are the basic unit of capacity in Kinesis. Make sure to create enough shards to accomodate your expected throughtput.


With On-Demand mode, the capcity of your stream updates dynamically depending on demand. The lower bound for writes is 4 MiB/second with an upper bound of 200MiB/second. You can request an increase of this quota up to 2 GiB/second by submitting a support request.

A single message (event) must not exceed 1MiB in size. This is a hard limit and you cannot request an increase.

see [AWS Kinesis documentation](https://docs.aws.amazon.com/streams/latest/dev/service-sizes-and-limits.html) for more information.
## Backpressure

You might see the `event_processing_duration_milliseconds` keeps on increasing and `batch_idle_in_channel_milliseconds` is in constant high value. In that case, Raccoon might get back-pressure from the publisher. If that happens, you can check the publisher, or you need to tune the publisher configuration on Raccoon.
