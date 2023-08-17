import unittest

from protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest, Event, SendEventResponse, Status, Code
from serde.json_serde import JsonSerde
from tests.unit.rest.client_test import get_static_uuid, get_static_time, get_stub_event_payload, get_stub_response


def get_event_request():
    request = SendEventRequest()
    request.req_guid = get_static_uuid()
    request.sent_time.FromNanoseconds(get_static_time())
    e = Event()
    e.type = "topic 1"
    e.event_bytes = b'{"random1": "abc", "xyz": 1}'
    request.events.append(e)
    return request


class JsonSerdeTest(unittest.TestCase):
    serde = JsonSerde()
    marshalled_event_request = """{
  "reqGuid": "17e2ac19-df8b-4a30-b111-fd7f5073d2f5",
  "sentTime": "2023-08-17T05:38:49.234986Z",
  "events": [
    {
      "eventBytes": "eyJyYW5kb20xIjogImFiYyIsICJ4eXoiOiAxfQ==",
      "type": "topic 1"
    }
  ]
}"""

    def test_serialise_of_input(self):
        event = {"random1": "abc", "xyz": 1}
        self.assertEqual(self.serde.serialise(event), b'{"random1": "abc", "xyz": 1}')

    def test_marshaling_of_proto_message(self):
        request = get_event_request()
        self.assertEqual(self.serde.marshal(request), self.marshalled_event_request)

    def test_unmarshaling_of_proto_message(self):
        stub_response = get_stub_response()._content
        unmarshalled_response = self.serde.unmarshal(stub_response, SendEventResponse())
        self.assertEqual(Status.STATUS_SUCCESS, unmarshalled_response.status)
        self.assertEqual(Code.CODE_OK, unmarshalled_response.code)
        self.assertEqual(get_static_time(), unmarshalled_response.sent_time)
        self.assertEqual(get_static_uuid(), unmarshalled_response.data["req_guid"])
