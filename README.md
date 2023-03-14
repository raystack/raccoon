# Raccoon

![build workflow](https://github.com/goto/raccoon/actions/workflows/build.yaml/badge.svg)
![package workflow](https://github.com/goto/raccoon/actions/workflows/package.yaml/badge.svg)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?logo=apache)](LICENSE)
[![Version](https://img.shields.io/github/v/release/goto/raccoon?logo=semantic-release)](Version)

Raccoon is high throughput, low-latency service that provides an API to ingest clickstream data from mobile apps, sites and publish it to Kafka. Raccoon uses the Websocket protocol for peer-to-peer communication and protobuf as the serialization format. It provides an event type agnostic API that accepts a batch (array) of events in protobuf format. Refer [here](https://github.com/goto/proton/tree/main/goto/raccoon) for proto definition format that Raccoon accepts.

<p align="center"><img src="./docs/static/assets/overview.svg" /></p>

## Key Features

- **Event Agnostic** - Raccoon API is event agnostic. This allows you to push any event with any schema.
- **Event Distribution** - Events are distributed to kafka topics based on the event meta-data
- **High performance** - Long running persistent, peer-to-peer connection reduce connection set up overheads. Websocket provides reduced battery consumption for mobile apps (based on usage statistics)
- **Guaranteed Event Delivery** - Server acknowledgements based on delivery. Currently it acknowledges failures/successes. Server can be augmented for zero-data loss or at-least-once guarantees.
- **Reduced payload sizes** - Protobuf based
- **Metrics:** - Built-in monitoring includes latency and active connections.

To know more, follow the detailed [documentation](https://goto.github.io/raccoon/)

## Use cases

Raccoon can be used as an event collector, event distributor and as a forwarder of events generated from mobile/web/IoT front ends as it provides an high volume, high throughput, low latency event-agnostic APIs. Raccoon can serve the needs of data ingestion in near-real-time. Some domains where Raccoon could be used is listed below

- Adtech streams: Where digital marketing data from external sources can be ingested into the organization backends
- Clickstream: Where user behavior data can be streamed in real-time
- Edge systems: Where devices (say in the IoT world) need to send data to the cloud.
- Event Sourcing: Such as Stock updates dashboards, autonomous/self-drive use cases

## Resources

Explore the following resources to get started with Raccoon:

- [Guides](https://goto.github.io/raccoon/guides/overview) provides guidance on deployment and client sample.
- [Concepts](https://goto.github.io/raccoon/concepts/architecture) describes all important Raccoon concepts.
- [Reference](https://goto.github.io/raccoon//reference/configurations) contains details about configurations, metrics and other aspects of Raccoon.
- [Contribute](https://goto.github.io/raccoon/contribute/contribution) contains resources for anyone who wants to contribute to Raccoon.

## Run with Docker

**Prerequisite**

- Docker installed

**Run Docker Image**

Raccoon provides Docker [image](https://hub.docker.com/r/goto/raccoon) as part of the release. Make sure you have Kafka running on your local and run the following.

```sh
# Download docker image from docker hub
$ docker pull gotocompany/raccoon

# Run the following docker command with minimal config.
$ docker run -p 8080:8080 \
  -e SERVER_WEBSOCKET_PORT=8080 \
  -e SERVER_WEBSOCKET_CONN_ID_HEADER=X-User-ID \
  -e PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS=host.docker.internal:9093 \
  -e EVENT_DISTRIBUTION_PUBLISHER_PATTERN=clickstream-%s-log \
  gotocompany/raccoon
```

**Run Docker Compose**
You can also use `docker-compose` on this repo. The `docker-compose` provides raccoon along with Kafka setup. Then, run the following command.

```sh
# Run raccoon along with kafka setup
$ make docker-run
# Stop the docker compose
$ make docker-stop
```

You can consume the published events from the host machine by using `localhost:9094` as kafka broker server. Mind the [topic routing](https://goto.github.io/raccoon/concepts/architecture#event-distribution) when you consume the events.

## Running locally

Prerequisite:

- You need to have [GO](https://golang.org/) 1.14 or above installed
- You need `protoc` [installed](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)

```sh
# Clone the repo
$ git clone https://github.com/goto/raccoon.git

# Build the executable
$ make

# Configure env variables
$ vim .env

# Run Raccoon
$ ./out/raccoon
```

**Note:** Read the detail of each configurations [here](https://goto.github.io/raccoon/reference/configuration).

## Running tests

```sh
# Running unit tests
$ make test

# Running integration tests
$ cp .env.test .env
$ make docker-run
$ INTEGTEST_BOOTSTRAP_SERVER=localhost:9094 INTEGTEST_HOST=localhost:8080 INTEGTEST_TOPIC_FORMAT="clickstream-%s-log" GRPC_SERVER_ADDR="localhost:8081" go test ./integration -v
```

## Contribute

Development of Raccoon happens in the open on GitHub, and we are grateful to the community for contributing bugfixes and improvements. Read below to learn how you can take part in improving Raccoon.

Read our [contributing guide](https://goto.github.io/raccoon/contribute/contribution) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to Raccoon.

To help you get your feet wet and get you familiar with our contribution process, we have a list of [good first issues](https://github.com/goto/raccoon/labels/good%20first%20issue) that contain bugs which have a relatively limited scope. This is a great place to get started.

This project exists thanks to all the [contributors](https://github.com/goto/raccoon/graphs/contributors).

## License

Raccoon is [Apache 2.0](LICENSE) licensed.
