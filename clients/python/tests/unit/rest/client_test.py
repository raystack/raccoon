import json
import time
import unittest

from unittest import mock
from unittest.mock import patch

import requests
from google.protobuf import timestamp_pb2

from raccoon_client import client
from raccoon_client.protos.raystack.raccoon.v1beta1.raccoon_pb2 import (
    SendEventRequest,
    Status,
    Code,
    SendEventResponse,
)
from raccoon_client.rest.client import RestClient
from raccoon_client.rest.option import RestClientConfigBuilder
from raccoon_client.serde.enum import Serialiser, WireType
from raccoon_client.serde.json_serde import JsonSerde
from raccoon_client.serde.protobuf_serde import ProtobufSerde


def get_marshalled_response():
    return b"\x08\x01\x10\x01\x18\x90\xe8\xc9\x97\xa8\xa2\x85\xbe\x17*0\n\x08req_guid\x12$17e2ac19-df8b-4a30-b111-fd7f5073d2f5"


def get_marshalled_request():
    return b"\n$17e2ac19-df8b-4a30-b111-fd7f5073d2f5\x12\x0b\x08\xe9\xe4\xf6\xa6\x06\x10\x90\xb4\x86p\x1a$\n\x14data bytes for click\x12\x0cclick-events"


def get_static_uuid():
    return "17e2ac19-df8b-4a30-b111-fd7f5073d2f5"


def get_static_time_ns():
    return 1692250729234986000


def get_static_time():
    return 1692276392


def get_stub_event_payload_json():
    return client.Event("random_topic", {"a": "abc"})


def get_stub_response_json():
    response = requests.Response()
    response.status_code = requests.status_codes.codes["ok"]
    json_response = {
        "status": 1,
        "code": 1,
        "sent_time": get_static_time(),
        "data": {"req_guid": get_static_uuid()},
    }
    json_string2 = json.dumps(json_response)
    response._content = json_string2
    return response


def get_stub_response_non_ok_json():
    response = requests.Response()
    response.status_code = requests.status_codes.codes["not_found"]
    json_response = {
        "status": Status.STATUS_ERROR,
        "code": Code.CODE_BAD_REQUEST,
        "sent_time": get_static_time(),
        "data": {"req_guid": get_static_uuid()},
    }
    json_string2 = json.dumps(json_response)
    response._content = json_string2
    return response


def get_stub_response_protobuf():
    response = requests.Response()
    response.status_code = requests.status_codes.codes["ok"]
    response._content = get_marshalled_response()
    return response


def get_stub_event_payload_protobuf():
    return client.Event(
        "random_topic",
        ProtobufSerde().unmarshal(
            get_marshalled_request(), SendEventRequest()
        ),  # sample proto serialised to bytes)
    )


class RestClientTest(unittest.TestCase):
    sample_url = "http://localhost:8080/api/v1/"
    max_retries = 3
    serialiser = Serialiser.JSON
    wire_type = WireType.JSON
    headers = {"X-Sample": "working"}

    def test_client_creation_success(self):
        client_config = (
            RestClientConfigBuilder()
            .with_url(self.sample_url)
            .with_serialiser(self.serialiser)
            .with_retry_count(self.max_retries)
            .with_wire_type(self.wire_type)
            .with_timeout(2.0)
            .with_headers(self.headers)
            .build()
        )
        rest_client = RestClient(client_config)
        self.assertEqual(
            rest_client.http_config.url, self.sample_url, "sample_urls do not match"
        )
        self.assertEqual(
            rest_client.session.adapters["https://"].max_retries.total,
            self.max_retries,
        )
        self.assertEqual(
            rest_client.session.adapters["http://"].max_retries.total, self.max_retries
        )
        self.assertEqual(
            type(rest_client.serde),
            self.serialiser.value,
            "serialiser is configured incorrectly",
        )
        self.assertEqual(
            type(rest_client.wire),
            self.wire_type.value,
            "wire type is configured incorrectly",
        )
        self.assertEqual(
            rest_client.http_config.timeout, 2.0, "timeout is configured incorrectly"
        )
        self.assertEqual(rest_client.http_config.headers, {"Content-Type": "application/json", "X-Sample": "working"})

    def test_client_creation_success_with_protobuf(self):
        client_config = (
            RestClientConfigBuilder()
            .with_url(self.sample_url)
            .with_serialiser(Serialiser.PROTOBUF)
            .with_retry_count(self.max_retries)
            .with_wire_type(WireType.PROTOBUF)
            .with_timeout(2.0)
            .with_headers({})
            .build()
        )
        rest_client = RestClient(client_config)
        self.assertEqual(
            rest_client.http_config.url, self.sample_url, "sample_urls do not match"
        )
        self.assertEqual(
            rest_client.session.adapters["https://"].max_retries.total,
            self.max_retries,
        )
        self.assertEqual(
            rest_client.session.adapters["http://"].max_retries.total, self.max_retries
        )
        self.assertEqual(
            type(rest_client.serde),
            Serialiser.PROTOBUF.value,
            "serialiser is configured incorrectly",
        )
        self.assertEqual(
            type(rest_client.wire),
            WireType.PROTOBUF.value,
            "wire type is configured incorrectly",
        )
        self.assertEqual(
            rest_client.http_config.timeout, 2.0, "timeout is configured incorrectly"
        )
        self.assertEqual(
            rest_client.http_config.headers, {"Content-Type": "application/proto"}
        )

    def test_client_creation_failure(self):
        builder = RestClientConfigBuilder().with_url(self.sample_url)
        self.assertRaises(ValueError, builder.with_serialiser, "JSON")
        self.assertRaises(ValueError, builder.with_wire_type, "PROTOBUF")
        self.assertRaises(ValueError, builder.with_retry_count, "five")
        self.assertRaises(ValueError, builder.with_timeout, 0.005)

    @patch("raccoon_client.rest.client.time.time_ns")
    def test_get_stub_request(self, time_ns):
        time_ns.return_value = get_static_time_ns()
        rest_client = self._get_rest_client()
        time_stamp = timestamp_pb2.Timestamp()  # pylint: disable=no-member
        time_stamp.FromNanoseconds(time.time_ns())
        with patch(
            "raccoon_client.rest.client.uuid.uuid4", return_value=get_static_uuid()
        ):
            req = rest_client._get_init_request()
            self.assertEqual(req.req_guid, get_static_uuid())
            self.assertEqual(req.sent_time.seconds, time_stamp.seconds)
            self.assertEqual(req.sent_time.nanos, time_stamp.nanos)

    def test_uniqueness_of_stub_request(self):
        rest_client = self._get_rest_client()
        req1 = rest_client._get_init_request()
        time.sleep(1)
        req2 = rest_client._get_init_request()
        self.assertNotEqual(req1.req_guid, req2.req_guid)
        self.assertNotEqual(req1.sent_time.nanos, req2.sent_time.nanos)
        self.assertNotEqual(req1.sent_time.seconds, req2.sent_time.seconds)

    def test_client_send_success_json(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.return_value = get_stub_response_json()
        event_arr = [get_stub_event_payload_json()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = get_static_uuid()
        req.req_guid = get_static_uuid()
        time_in_ns = time.time_ns()
        req.sent_time.FromNanoseconds(time_in_ns)
        expected_req.sent_time.FromNanoseconds(time_in_ns)
        with patch(
            "raccoon_client.rest.client.requests.session", return_value=session_mock
        ):
            rest_client = self._get_rest_client()
            expected_req.events.append(
                rest_client._convert_to_event_pb(get_stub_event_payload_json())
            )
            serialised_data = JsonSerde().marshal(expected_req)
            rest_client._get_init_request = mock.MagicMock()
            rest_client._get_init_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            rest_client._parse_response.return_value = [SendEventResponse(), None]
            rest_client.send(event_arr)
            post.assert_called_once_with(
                url=self.sample_url,
                data=serialised_data,
                headers={"Content-Type": "application/json"},
                timeout=2.0,
            )
            rest_client._parse_response.assert_called_once_with(post.return_value)

    def test_client_send_success_protobuf(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.return_value = get_stub_response_protobuf()
        event_arr = [get_stub_event_payload_protobuf()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = get_static_uuid()
        req.req_guid = get_static_uuid()
        time_in_ns = time.time_ns()
        req.sent_time.FromNanoseconds(time_in_ns)
        expected_req.sent_time.FromNanoseconds(time_in_ns)
        with patch(
            "raccoon_client.rest.client.requests.session", return_value=session_mock
        ):
            rest_client = self._get_rest_client(
                serialiser=Serialiser.PROTOBUF, wire_type=WireType.PROTOBUF
            )
            expected_req.events.append(
                rest_client._convert_to_event_pb(get_stub_event_payload_protobuf())
            )
            serialised_data = expected_req.SerializeToString()
            rest_client._get_init_request = mock.MagicMock()
            rest_client._get_init_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            rest_client._parse_response.return_value = [SendEventResponse(), None]
            rest_client.send(event_arr)
            post.assert_called_once_with(
                url=self.sample_url,
                data=serialised_data,
                headers={"Content-Type": "application/proto"},
                timeout=2.0,
            )
            rest_client._parse_response.assert_called_once_with(post.return_value)

    def test_client_send_success_json_serialiser_protobuf_wire(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.return_value = get_stub_response_protobuf()
        event_arr = [get_stub_event_payload_json()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = get_static_uuid()
        req.req_guid = get_static_uuid()
        time_in_ns = time.time_ns()
        req.sent_time.FromNanoseconds(time_in_ns)
        expected_req.sent_time.FromNanoseconds(time_in_ns)
        with patch(
            "raccoon_client.rest.client.requests.session", return_value=session_mock
        ):
            rest_client = self._get_rest_client(
                serialiser=Serialiser.JSON, wire_type=WireType.PROTOBUF
            )
            expected_req.events.append(
                rest_client._convert_to_event_pb(get_stub_event_payload_json())
            )
            serialised_data = expected_req.SerializeToString()
            rest_client._get_init_request = mock.MagicMock()
            rest_client._get_init_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            rest_client._parse_response.return_value = [SendEventResponse(), None]
            rest_client.send(event_arr)
            post.assert_called_once_with(
                url=self.sample_url,
                data=serialised_data,
                headers={"Content-Type": "application/proto"},
                timeout=2.0,
            )
            rest_client._parse_response.assert_called_once_with(post.return_value)

    def test_client_send_success_protobuf_serialiser_json_wire(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.return_value = get_stub_response_json()
        event_arr = [get_stub_event_payload_protobuf()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = get_static_uuid()
        req.req_guid = get_static_uuid()
        time_in_ns = time.time_ns()
        req.sent_time.FromNanoseconds(time_in_ns)
        expected_req.sent_time.FromNanoseconds(time_in_ns)
        with patch(
            "raccoon_client.rest.client.requests.session", return_value=session_mock
        ):
            rest_client = self._get_rest_client(
                serialiser=Serialiser.PROTOBUF, wire_type=WireType.JSON
            )
            expected_req.events.append(
                rest_client._convert_to_event_pb(get_stub_event_payload_protobuf())
            )
            serialised_data = JsonSerde().marshal(expected_req)
            rest_client._get_init_request = mock.MagicMock()
            rest_client._get_init_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            rest_client._parse_response.return_value = [SendEventResponse(), None]
            rest_client.send(event_arr)
            post.assert_called_once_with(
                url=self.sample_url,
                data=serialised_data,
                headers={"Content-Type": "application/json"},
                timeout=2.0,
            )
            rest_client._parse_response.assert_called_once_with(post.return_value)

    def test_client_send_connection_failure(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.side_effect = ConnectionError("error connecting to host")
        event_arr = [get_stub_event_payload_json()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = get_static_uuid()
        req.req_guid = get_static_uuid()
        time_in_ns = time.time_ns()
        req.sent_time.FromNanoseconds(time_in_ns)
        expected_req.sent_time.FromNanoseconds(time_in_ns)
        with patch(
            "raccoon_client.rest.client.requests.session", return_value=session_mock
        ):
            rest_client = self._get_rest_client()
            expected_req.events.append(
                rest_client._convert_to_event_pb(get_stub_event_payload_json())
            )
            serialised_data = JsonSerde().marshal(expected_req)
            rest_client._get_init_request = mock.MagicMock()
            rest_client._get_init_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            self.assertRaises(ConnectionError, rest_client.send, event_arr)
            post.assert_called_once_with(
                url=self.sample_url,
                data=serialised_data,
                headers={"Content-Type": "application/json"},
                timeout=2.0,
            )
            rest_client._parse_response.assert_not_called()

    def test_parse_response_json(self):
        resp = get_stub_response_json()
        rest_client = self._get_rest_client()
        deserialised_response, err = rest_client._parse_response(resp)
        self.assertEqual(deserialised_response.status, Status.STATUS_SUCCESS)
        self.assertEqual(deserialised_response.data["req_guid"], get_static_uuid())
        self.assertEqual(deserialised_response.sent_time, get_static_time())
        self.assertEqual(deserialised_response.code, Code.CODE_OK)
        self.assertIsNone(err)

    def test_parse_response_protobuf(self):
        resp = get_stub_response_protobuf()
        rest_client = self._get_rest_client(wire_type=WireType.PROTOBUF)
        deserialised_response, err = rest_client._parse_response(resp)
        self.assertEqual(deserialised_response.status, Status.STATUS_SUCCESS)
        self.assertEqual(deserialised_response.data["req_guid"], get_static_uuid())
        self.assertEqual(deserialised_response.sent_time, get_static_time_ns())
        self.assertEqual(deserialised_response.code, Code.CODE_OK)
        self.assertIsNone(err)

    def test_parse_response_for_non_ok_status(self):
        resp = get_stub_response_non_ok_json()
        rest_client = self._get_rest_client()
        deserialised_response, err = rest_client._parse_response(resp)
        self.assertEqual(deserialised_response.status, Status.STATUS_ERROR)
        self.assertEqual(deserialised_response.data["req_guid"], get_static_uuid())
        self.assertEqual(deserialised_response.sent_time, get_static_time())
        self.assertEqual(deserialised_response.code, Code.CODE_BAD_REQUEST)
        self.assertIsNotNone(err)
        self.assertEqual(err.status_code, 404)

    def _get_rest_client(self, serialiser=Serialiser.JSON, wire_type=WireType.JSON):
        client_config = (
            RestClientConfigBuilder()
            .with_url(self.sample_url)
            .with_serialiser(serialiser)
            .with_retry_count(self.max_retries)
            .with_wire_type(wire_type)
            .with_timeout(2.0)
            .build()
        )
        return RestClient(client_config)
