# Troubleshooting

## Scale Up Racoon

Internally, Raccoon has 3 main components that affect the capacity. The server, worker, and publisher. Each component has configurations that can be tune if necessary. Since those 3 components are forming a pipe, you need to make sure none of the components become a bottleneck.

To know the right configuration, you need to simulate with similar throughput as production. You can tune the configuration accordingly.

Following are details of what you can tune.

### Server

Raccoon is using WebSocket as a communication protocol from client to server. Websocket requires maintaining long-running connections. Each connection costs the OS an open file descriptor. When you reach the limit of the configured open file descriptor, the server won't be able to accept a new connection. By default, OS limit the number of the open file descriptor. You can look up how to increase the max open file descriptor. On Unix, you can do `ulimit -n` to check max open file descriptor and `ulimit -n <number>` to set a new limit.

Apart from OS configuration, there are configurations you can tune on Raccoon:

* [SERVER\_WEBSOCKET\_MAX\_CONN](https://goto.gitbook.io/raccoon/reference/configurations#server_websocket_max_conn) To limit Raccoon resource utilization, we enforce a limit on WebSocket connection. The default value is 30000; adjust it if necessary.

### Worker

After the request is deserialized, the server puts the events on the buffer channel. The worker process events from the channel and publish them to Kafka. You can think of the worker and the channel as a buffer in case the publisher slows down temporarily.

* [WORKER\_BUFFER\_CHANNEL\_SIZE](https://goto.gitbook.io/raccoon/reference/configurations#worker_buffer_channel_size) Buffer before the events get processed. The more the size, the longer it can tolerate a temporary spike or slow down.
* [WORKER\_POOL\_SIZE](https://goto.gitbook.io/raccoon/reference/configurations#worker_pool_size) The worker will call the publisher client and wait synchronously. Increase this according to the throughput.

### Publisher

Currently, Raccoon is using [Librd Kafka client Go wrapper](https://github.com/confluentinc/confluent-kafka-go) as publisher client. There is plenty of guides out there to tune Kafka producer. Here are some configurations you can tune.

* [PUBLISHER\_KAFKA\_CLIENT\_BATCH\_NUM\_MESSAGES](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md)
* [KAFKA\_CLIENT\_ACKS](https://goto.gitbook.io/raccoon/reference/configurations#publisher_kafka_client_acks)
* KAFKA\_CLIENT\_LINGER\_MS
* [PUBLISHER_KAFKA\_CLIENT_\*](https://goto.gitbook.io/raccoon/reference/configurations#publisher_kafka_client_) You can put any [librd kafka configuration](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md) by replacing \* and change the delimiter to underscore.

## Backpressure

You might see the `event_processing_duration_milliseconds` keeps on increasing and `batch_idle_in_channel_milliseconds` is in constant high value. In that case, Raccoon might get backpressure from the publisher. If that happens, you can check the publisher, or you need to tune the publisher configuration on Raccoon.

