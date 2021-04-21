# Introduction

Raccoon is high throughput, low-latency service that provides an API to ingest clickstream data from mobile apps, sites and publish it to Kafka. Raccoon uses the Websocket protocol for peer-to-peer communication and protobuf as the serialization format. It provides an event type agnostic API that accepts a batch (array) of events in protobuf format. Refer [here](https://github.com/odpf/proton/tree/main/odpf/raccoon) for proto definition format that Raccoon accepts.

<p align="center"><img src="./docs/assets/overview.png" /></p>

## Key Features

* **Event Agnostic:** Raccoon API is event agnostic. This allows you to push any event with any schema.
* **Metrics:** Built in monitoring includes latency and active connections.

To know more, follow the detailed [documentation](docs) 

## Usage

Explore the following resources to get started with Raccoon:

* [Guides](docs/guides) provides guidance on deployment and client sample.
* [Concepts](docs/concepts) describes all important Raccoon concepts.
* [Reference](docs/reference) contains details about configurations, metrics and other aspects of Raccoon.
* [Contribute](docs/contribute/contribution.md) contains resources for anyone who wants to contribute to Raccoon.