import unittest

from raccoon_client.protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest, Event, SendEventResponse, \
    Status, Code
from raccoon_client.serde.protobuf_serde import ProtobufSerde
from tests.unit.rest.client_test import get_static_uuid, get_static_time_ns


def get_stub_request() -> SendEventRequest:
    req = SendEventRequest()
    req.req_guid = get_static_uuid()
    req.sent_time.FromNanoseconds(get_static_time_ns())
    e = Event()
    e.type = "click-events"
    e.event_bytes = bytes("data bytes for click", "utf-8")
    req.events.append(e)
    return req


def get_marshalled_response():
    return b'\x08\x01\x10\x01\x18\x90\xe8\xc9\x97\xa8\xa2\x85\xbe\x17*0\n\x08req_guid\x12$17e2ac19-df8b-4a30-b111-fd7f5073d2f5'


def get_marshalled_request():
    return b'\n$17e2ac19-df8b-4a30-b111-fd7f5073d2f5\x12\x0b\x08\xe9\xe4\xf6\xa6\x06\x10\x90\xb4\x86p\x1a$\n\x14data bytes for click\x12\x0cclick-events'


class ProtobufSerdeTest(unittest.TestCase):
    serde = ProtobufSerde()

    def test_serialisation_of_proto_input(self):
        event = get_stub_request()
        serialised_data = self.serde.serialise(event)
        expected_serialised_data = get_marshalled_request()
        self.assertEqual(expected_serialised_data, serialised_data)

    def test_serialisation_non_proto_input(self):
        event = {"a": "b"}
        self.assertRaises(ValueError, self.serde.serialise, event)

    def test_marshalling_of_payload(self):
        event = get_stub_request()
        marshalled_data = self.serde.marshal(event)
        expected_marshalled_data = get_marshalled_request()
        self.assertEqual(expected_marshalled_data, marshalled_data)

    def test_unmarshalling_of_payload(self):
        x = [8, 1, 16, 1, 24, 204, 153, 179, 102, 42, 48, 10, 8, 114, 101, 113, 95, 103, 117, 105, 100, 18, 36, 54, 49,
             54, 57, 49, 98, 50, 51, 45, 52, 98, 55, 102, 45, 52, 54, 99, 102, 45, 97, 102, 57, 51, 45, 97, 98, 97, 56,
             55, 52, 99, 50, 52, 49, 56, 54]
        marshalled_response = bytes(x)
        unmarshalled_response = self.serde.unmarshal(marshalled_response, SendEventResponse())
        self.assertEqual(Status.STATUS_SUCCESS, unmarshalled_response.status)
        self.assertEqual(Code.CODE_OK, unmarshalled_response.code)
        self.assertEqual(get_static_time_ns(), unmarshalled_response.sent_time)
        self.assertEqual(get_static_uuid(), unmarshalled_response.data["req_guid"])

    def test_correct_content_type(self):
        self.assertEqual("application/proto", self.serde.get_content_type())
