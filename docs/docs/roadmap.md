# Roadmap

In the following section, you can learn what features we're working on, what stage they're in, and when we expect to bring them to you. Have any questions or comments about items on the roadmap? Join the [discussions](https://github.com/raystack/raccoon/discussions) on the Raccoon Github forum.

We’re planning to iterate on the format of the roadmap itself, and we see the potential to engage more in discussions about the future of Raccoon features. If you have feedback about this roadmap section itself, such as how the issues are presented, let us know through [discussions](https://github.com/raystack/raccoon/discussions).

## Vision

We want to enable Raccoon as the preferred event collector, event distributor that provides high volume, high throughput, low latency protocol-agnostic, event-agnostic APIs for data ingestion in near-real-time. With this vision, Raccoon can serve the needs of Adtech streams - Where digital marketing data from external sources can be ingested into the organization backends Clickstream - Where user behavior data can be streamed in real-time Edge networks - Where devices \(say in the IoT world\) need to send data to the cloud. Event Sourcing systems - Such as Stock updates dashboards, autonomous/self-drive use cases

![](/assets/raccoon_vision.png)

### Raccoon 1.x

- Support for HTTP, gRPC
- Support for json, protobuf formats
- Extendable event distribution
- Extendable event filtering capability
- Enable Raccoon to replay lost events with zero-data-loss capability
- Adopt Raccoon to publish to different transport systems
- Enables Raccoon to provide extendable data formatters. eg. JSON to proto
