# Structure

This document describes high-level code structure of the project. You'll find this part useful when you want to get started to contribute to Raccoon.

## Highlevel View

The core structure of Raccoon is the server itself. After the server is started, data flows from `websocket` to `worker` to `publisher`. `websocket` manages websocket server, handle incoming connection, and incoming request. `worker` acts as a buffer and interface for various types of server and publisher down the roadmap. `publisher` contains logic to publish the events to the downstream pipeline. ![high-level](/assets/structure.svg) All the components above are initialized on `app`. `app` package is the starting point of Raccoon.

## Code Map

This section talks briefly about the content of various important packages. You can use this to guess the code location when you need to make changes.

### `http`

Contains all the http related code including code related to `websocket`, `rest` and `grpc`. It also has code pertaining to the http server that serves both WebSocket and REST APIs.

#### `http/websocket`

Contains server-related code along with request/response handlers and [connection management](architecture.md#connections).

#### `http/rest`

Contains server-side code along with resquest/response handler for the REST endpoint.

#### `http/gRPC`

Contains server-side handlers for gRPC server.

### `worker`

Buffer from when the events are processed and before events are published. This will also act as interface that connects server and publisher when in the future. Currently, `worker` is tightly coupled with `websocket` server and `kafka` publisher.

### `publisher`

Does the actual publishing to the downstream pipeline. Currently, only support Kafka publisher.

### `app`

The starting point of Raccoon. It initializes server, worker, publisher, and other components that require initialization like a metric reporter.

### `config`

Load and store application configurations. Contains mapping of environment configuration with configuration on the code.

### `serialization`

Contains the common serialization code for both JSON and Protobufs along with common interface.

### `deserialization`

Contains the common deserialization code along with common interface.

### `identification`

Contains the code for connection identification.

## Code Generation

### Request, Response, and Events Proto

Raccoon depends on [Proton](https://github.com/goto/proton/tree/main/goto/raccoon) repository. Proton is a repository to store all gotocompany Protobuf files. Code to serde the request and response are generated using Protobuf. You can check how the code is generated on `Makefile`.
