# Overview

## Introduction
Raccoon provides a suite of client libraries designed to help developers easily send clickstream events to its low-latency, high-throughput event ingestion service. Whether you’re building real-time analytics, tracking user behavior, or processing large-scale event data, Raccoon's clients offer flexible and efficient integration via WebSocket, REST, and gRPC APIs.

## Key Features

- **Multi-Protocol Support**: WebSocket, REST, and gRPC are available in all clients, allowing you to choose the best fit for your application’s needs.
- **Ease of Integration**: Designed with simplicity in mind, the clients integrate easily into existing projects with minimal configuration.
- **Reliability**: Each client includes retry mechanisms and error handling to ensure events are delivered reliably, even in the face of transient failures.

## Getting Started

To start using Raccoon's client libraries, check out the detailed installation and usage instructions for each supported language:

- [Golang](clients/golang.md)
- [Python](clients/python.md)
- [Java](clients/java.md)
- [JavaScript](clients/javascript.md)

By leveraging Raccoon’s clients, you can focus on building your applications while Raccoon efficiently handles the ingestion of your clickstream events.
