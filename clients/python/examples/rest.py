from raccoon_client.client import Event
from raccoon_client.rest.client import RestClient
from raccoon_client.rest.option import RestClientConfigBuilder
from raccoon_client.serde.enum import Serialiser, WireType

if __name__ == '__main__':
    event_data = {"a": "field a", "b": "field b"}

    config = RestClientConfigBuilder().with_url("http://localhost:8080/api/v1/events")\
        .with_serialiser(Serialiser.PROTOBUF)\
        .with_wire_type(WireType.JSON)\
        .build()  # other parameters supported by the config builder can be checked in its method definition.
    rest_client = RestClient(config)
    topic_to_publish_to = "test_topic_2"
    e = Event(topic_to_publish_to, event_data)
    req_id, response, raw = rest_client.send([e])
