import unittest

from raccoon_client.protos.raystack.raccoon.v1beta1.raccoon_pb2 import (
    SendEventRequest,
    Event,
    SendEventResponse,
    Status,
    Code,
)
from raccoon_client.serde.json_serde import JsonSerde
from tests.unit.rest.client_test import (
    get_static_uuid,
    get_static_time_ns,
    get_stub_response_json,
    get_static_time,
)


def get_event_request():
    request = SendEventRequest()
    request.req_guid = get_static_uuid()
    request.sent_time.FromNanoseconds(get_static_time_ns())
    event = Event()
    event.type = "topic 1"
    event.event_bytes = b'{"random1": "abc", "xyz": 1}'
    request.events.append(event)
    return request


class JsonSerdeTest(unittest.TestCase):
    serde = JsonSerde()
    marshalled_event_request = """{"req_guid": "17e2ac19-df8b-4a30-b111-fd7f5073d2f5", "sent_time": {"seconds": 1692250729, "nanos": 234986000}, "events": [{"event_bytes": "eyJyYW5kb20xIjogImFiYyIsICJ4eXoiOiAxfQ==", "type": "topic 1"}]}"""
    serialised_event_request = b'{\n  "reqGuid": "17e2ac19-df8b-4a30-b111-fd7f5073d2f5",\n  "sentTime": "2023-08-17T05:38:49.234986Z",\n  "events": [\n    {\n      "eventBytes": "eyJyYW5kb20xIjogImFiYyIsICJ4eXoiOiAxfQ==",\n      "type": "topic 1"\n    }\n  ]\n}'

    def test_serialise_of_input(self):
        event = {"random1": "abc", "xyz": 1}
        self.assertEqual(self.serde.serialise(event), b'{"random1": "abc", "xyz": 1}')

    def test_serialise_of_proto_object(self):
        event = get_event_request()
        serialised_proto = self.serde.serialise(event)
        self.assertEqual(serialised_proto, self.serialised_event_request)

    def test_marshaling_of_proto_message(self):
        request = get_event_request()
        self.assertEqual(self.serde.marshal(request), self.marshalled_event_request)

    def test_unmarshalling_into_proto_message(self):
        stub_response = get_stub_response_json()._content
        unmarshalled_response = self.serde.unmarshal(stub_response, SendEventResponse())
        self.assertEqual(Status.STATUS_SUCCESS, unmarshalled_response.status)
        self.assertEqual(Code.CODE_OK, unmarshalled_response.code)
        self.assertEqual(get_static_time(), unmarshalled_response.sent_time)
        self.assertEqual(get_static_uuid(), unmarshalled_response.data["req_guid"])

    def test_content_type_for_json(self):
        self.assertEqual("application/json", self.serde.get_content_type())
