# Introduction

Raccoon is high throughput, low-latency service that provides an API to ingest clickstream data from mobile apps, sites and publish it to Kafka. Raccoon uses the Websocket protocol for peer-to-peer communication and protobuf as the serialization format. It provides an event type agnostic API that accepts a batch (array) of events in protobuf format. Refer [here](https://github.com/odpf/proton/tree/main/odpf/raccoon) for proto definition format that Raccoon accepts.

![Overiew](./assets/overview.svg)

## Key Features

* **Event Agnostic** - Raccoon API is event agnostic. This allows you to push any event with any schema.
* **Event Distribution** - Events are distributed to kafka topics based on the event meta-data
* **High performance** - Long running persistent, peer-to-peer connection reduce connection set up overheads. Websocket provides reduced battery consumption for mobile apps (based on usage statistics)
* **Guaranteed Event Delivery** - Server acknowledgements based on delivery. Currently it acknowledges failures/successes. Server can be augmented for zero-data loss or at-least-once guarantees.
* **Reduced payload sizes** - Protobuf based
* **Metrics:** - Built-in monitoring includes latency and active connections.

## Use Cases
Raccoon can be used as an event collector, event distributor and as a forwarder of events generated from mobile/web/IoT front ends as it provides an high volume, high throughput, low latency event-agnostic APIs. Raccoon can serve the needs of data ingestion in near-real-time. Some domains where Raccoon could be used is listed below

* Adtech streams: Where digital marketing data from external sources can be ingested into the organization backends 
* Clickstream: Where user behavior data can be streamed in real-time 
* Edge systems: Where devices (say in the IoT world) need to send data to the cloud. 
* Event Sourcing: Such as Stock updates dashboards, autonomous/self-drive use cases

## Usage

Explore the following resources to get started with Raccoon:

* [Guides](docs/guides) provides guidance on deployment and client sample.
* [Concepts](docs/concepts) describes all important Raccoon concepts.
* [Reference](docs/reference) contains details about configurations, metrics and other aspects of Raccoon.
* [Contribute](docs/contribute/contribution.md) contains resources for anyone who wants to contribute to Raccoon.