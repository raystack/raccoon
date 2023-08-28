from raccoon_client.client import Event
from raccoon_client.protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest
from raccoon_client.rest.client import RestClient
from raccoon_client.rest.option import RestClientConfigBuilder
from raccoon_client.serde.enum import Serialiser, WireType


def example_json_serialiser_json_wire():
    event_data = {"a": "field a", "b": "field b"}

    config = (
        RestClientConfigBuilder()
        .with_url("http://localhost:8080/api/v1/events")
        .with_serialiser(Serialiser.JSON)
        .with_wire_type(WireType.JSON)
        .build()
    )  # other parameters supported by the config builder can be checked in its method definition.
    rest_client = RestClient(config)
    topic_to_publish_to = "test_topic_2"
    e = Event(topic_to_publish_to, event_data)
    req_id, response, raw = rest_client.send([e])
    return req_id, response, raw


def example_protobuf_serialiser_protobuf_wire():
    event_data = (
        SendEventRequest()
    )  # sample generated proto class which is an event to send to raccoon
    event_data.sent_time = 1000
    event_data.req_guid = "some string"

    config = (
        RestClientConfigBuilder()
        .with_url("http://localhost:8080/api/v1/events")
        .with_serialiser(Serialiser.PROTOBUF)
        .with_wire_type(WireType.PROTOBUF)
        .with_timeout(10.0)
        .with_retry_count(3)
        .with_headers({"Authorization", "TOKEN"})
        .build()
    )  # other parameters supported by the config builder can be checked in its method definition.
    rest_client = RestClient(config)
    topic_to_publish_to = "test_topic_2"
    e = Event(topic_to_publish_to, event_data)
    req_id, response, raw = rest_client.send([e])
    return req_id, response, raw


def example_protobuf_serialiser_json_wire():
    event_data = (
        SendEventRequest()
    )  # sample generated proto class which is an event to send to raccoon
    event_data.sent_time = 1000
    event_data.req_guid = "some string"

    config = (
        RestClientConfigBuilder()
        .with_url("http://localhost:8080/api/v1/events")
        .with_serialiser(Serialiser.PROTOBUF)
        .with_wire_type(WireType.JSON)
        .build()
    )  # other parameters supported by the config builder can be checked in its method definition.
    rest_client = RestClient(config)
    topic_to_publish_to = "test_topic_2"
    e = Event(topic_to_publish_to, event_data)
    req_id, response, raw = rest_client.send([e])
    return req_id, response, raw


if __name__ == "__main__":
    example_json_serialiser_json_wire()
    example_protobuf_serialiser_protobuf_wire()
    example_protobuf_serialiser_json_wire()
