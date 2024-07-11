---
id: introduction
slug: /
---

# Introduction

Raccoon is a high-throughput, low-latency service that provides an API to ingest clickstream data from mobile apps and sites and publish it to Kafka. Raccoon uses the WebSocket protocol for peer-to-peer communication and protobuf as the serialization format. It provides an event-type agnostic API that accepts a batch \(array\) of events in protobuf format. Refer [here](https://github.com/raystack/proton/tree/main/raystack/raccoon) for the protobuf definition format that Raccoon accepts.

![Overiew](/assets/overview.svg)

## Key Features

- **Event Agnostic** - Raccoon API is event agnostic. This allows you to push any event with any schema.
- **Event Distribution** - Events are distributed to Kafka topics based on the event meta-data.
- **High performance** - Long-running persistent, peer-to-peer connection reduce connection set up overheads. Websocket provides reduced battery consumption for mobile apps \(based on usage statistics\)
- **Guaranteed Event Delivery** - Server acknowledgments based on delivery. Currently, it acknowledges failures/successes. Additionally, users can augment the server for zero-data loss or at-least-once guarantees.
- **Reduced payload sizes** - Protobuf based
- **Metrics:** - Built-in monitoring includes latency and active connections.

## Use Cases

Raccoon can be used as an event collector, event distributor, and forwarder of events generated from mobile/web/IoT front-ends as it provides a high volume, high throughput, low latency event-agnostic APIs. In addition, it can serve the needs of data ingestion in near-real-time. Some domains where Raccoon could be used are listed below.

- Adtech streams: Where users can ingest digital marketing data from external sources into the organization's backends
- Clickstream: Where apps can stream user behavior data in real-time
- Edge systems: Where devices \(say in the IoT world\) need to send data to the cloud.
- Event Sourcing: Such as stock update dashboards, autonomous/self-drive use cases

## Usage

Explore the following resources to get started with Raccoon:

- [Guides](guides/overview.md) provide information on deployment and client samples.
- [Concepts](concepts/architecture.md) describe all important Raccoon concepts.
- [Reference](reference/configurations.md) contains details about configurations, metrics, and other aspects of Raccoon.
- [Contribute](contribute/contribution.md) contains resources for anyone who wants to contribute to Raccoon.
