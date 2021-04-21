# Introduction

Raccoon is high throughput, low-latency service that provides an API to ingest clickstream data from mobile apps, sites and publish it to Kafka. Raccoon uses the Websocket protocol for peer-to-peer communication and protobuf as the serialization format. It provides an event type agnostic API that accepts a batch \(array\) of events in protobuf format. Refer [here](https://github.com/odpf/proton/tree/main/odpf/raccoon) for proto definition format that Raccoon accepts.

## Key Features

* **Event Agnostic:** Raccoon API is event agnostic. This allows you to push any event with any schema.
* **Metrics:** Built in monitoring includes latency and active connections.

To know more, follow the detailed [documentation](https://github.com/odpf/raccoon/tree/081b02c61ad669301379b304bb0ff839ca44d02c/docs/docs/README.md)

## Usage

Explore the following resources to get started with Raccoon:

* [Guides](https://github.com/odpf/raccoon/tree/081b02c61ad669301379b304bb0ff839ca44d02c/docs/docs/guides/README.md) provides guidance on deployment and client sample.
* [Concepts](https://github.com/odpf/raccoon/tree/081b02c61ad669301379b304bb0ff839ca44d02c/docs/docs/concepts/README.md) describes all important Raccoon concepts.
* [Reference](https://github.com/odpf/raccoon/tree/081b02c61ad669301379b304bb0ff839ca44d02c/docs/docs/reference/README.md) contains details about configurations, metrics and other aspects of Raccoon.
* [Contribute](https://github.com/odpf/raccoon/tree/081b02c61ad669301379b304bb0ff839ca44d02c/docs/docs/contribute/contribution.md) contains resources for anyone who wants to contribute to Raccoon.

