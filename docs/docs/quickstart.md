# Quickstart

This document will guide you to get Raccoon along with Kafka setup running locally. This document assumes that you have installed Docker and Kafka with `host.docker.internal` [advertised](https://www.confluent.io/blog/kafka-listeners-explained/) on your machine.

## Run Raccoon With Docker

Run the following command. Make sure to set `PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS` according to your local Kafka setup.

```bash
$ docker run -p 8080:8080 \
  -e SERVER_WEBSOCKET_CONN_ID_HEADER=X-User-ID \
  -e PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS=host.docker.internal:9092 \
  -e EVENT_DISTRIBUTION_PUBLISHER_PATTERN=clickstream-log \
  gotocompany/raccoon:latest
```

To test whether the service is running or not, you can try to ping the server.

```bash
$ curl http://localhost:8080/ping
```

## Publishing Your First Event

Currently, Raccoon doesn't come with a library client. To start publishing events to Raccoon, we provide you an [example of a go client](https://github.com/goto/raccoon/tree/main/docs/example) that you can refer to. You can also run the example right away if you have Go installed on your machine.

```bash
# `cd` on the client example directory and run the following
$ go run main.go sample.pb.go
```

To verify the event published by Raccoon. First, you need to start a Kafka listener.

```bash
$ kafka-console-consumer --bootstrap-server localhost:9092 --topic clickstream-log
```

## Where To Go Next

For more detail about publishing events to Raccoon, you can read the [detailed document](https://goto.gitbook.io/raccoon/guides/publishing) under the guides section. To understand more about how Raccoon work, you can go to the [architecture document](https://goto.gitbook.io/raccoon/concepts/architecture).

