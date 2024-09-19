# Python

## Requirements
Make sure that Python version `>=3.9` is installed on your system. See [installation instructions](https://docs.python.org/3.9/using/unix.html#getting-and-installing-the-latest-version-of-python) on Python's website for more information.

## Installation
Install Raccoon's Python client [pip](https://docs.python.org/3/installing/index.html)
```bash
$ pip install raccoon_client
```

Or if you prefer using [poetry](https://python-poetry.org/docs/)
```bash
$ poetry add raccoon_client
```
## Usage

### Quickstart

Below is a self contained example of Raccoon's Python client that uses the REST API to publish events

```python title="quickstart.py"
from raccoon_client.client import Event
from raccoon_client.protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest
from raccoon_client.rest.client import RestClient
from raccoon_client.rest.option import RestClientConfigBuilder
from raccoon_client.serde.enum import Serialiser, WireType

def run():
    config = (
        RestClientConfigBuilder()
        .with_url("http://localhost:8080/api/v1/events")
        .with_serialiser(Serialiser.JSON)
        .with_wire_type(WireType.JSON)
        .build()
    )  
    client = RestClient(config)
    data = [
        {
            "a": "field a",
            "b": "field b",
        }
    ]
    topic = "test_topic_2"
    events = [Event(topic, event)]
    req_id, response, raw = client.send(events)
    print(req_id, response, raw)

if __name__ == "__main__":
    run()
```

### Guide

#### Creating a client

Raccoon's Python only supports sending events over Raccoon's HTTP/JSON (REST) API.

To create a client, you must pass the `RestClientConfig` object to the client constructor `RestClient`.

To build the client config, use `RestClientConfigBuilder`. `RestClientConfigBuilder` uses a [builder pattern](https://en.wikipedia.org/wiki/Builder_pattern) along with a [fluent interface API](https://en.wikipedia.org/wiki/Fluent_interface) to help build the config object. Here's a minimal example:

```python
from raccoon_client.rest.client import RestClient
from raccoon_client.rest.option import RestClientConfigBuilder
from raccoon_client.serde.enum import Serialiser, WireType

config = (
    RestClientConfigBuilder()
        .with_url("http://localhost:8080/api/v1/events")
        .with_serializer(Serializer.JSON)
        .with_wire_type(WireType.JSON)
        .build()
)
client = RestClient(config)
```
#### Publishing events

To publish events, create a list of `raccoon_client.client.Event` and pass it to `RestClient.send()` method. Each event has a `type` property and an `event` property.

`type` denotes the event type. This is used by raccoon to route the event to a specific topic downstream. `event` field contains the payload or raw event data. This data is serialised by the `serializer` that's configured on the client. 

The following table lists which serializer to use for a given payload type.

| Message Type | Serializer |
| --- | --- |
| JSON | `raccoon_client.serde.enum.Serialiser.JSON` |
| Protobuf | `raccoon_client.serde.emum.Serialiser.PROTOBUF`|

Once a client is constructed with a specific kind of serializer, you may only pass it events of that specific type. In particular, for `JSON` serialiser the event data must be a python dict. While for `PROTOBUF` serialiser the event data must be a protobuf message.

### Examples
You can find examples for Raccoon's python client [here](https://github.com/raystack/raccoon/blob/main/clients/python/examples/rest.py)