import json
import time
import unittest

from unittest import mock
from unittest.mock import patch

import requests
from google.protobuf import json_format, timestamp_pb2

import client
from protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest, SendEventResponse, Status
from rest.client import RestClient
from rest.option import RestClientConfigBuilder
from serde.enum import Serialiser, WireType


def _get_static_uuid():
    return "17e2ac19-df8b-4a30-b111-fd7f5073d2f5"


def _get_static_time():
    return 1692182259


class RestClientTest(unittest.TestCase):
    sample_url = "http://localhost:8080/api/v1/"
    max_retries = 3
    serialiser = Serialiser.JSON
    wire_type = WireType.JSON

    def test_client_creation_success(self):
        client_config = RestClientConfigBuilder(). \
            with_url(self.sample_url). \
            with_serialiser(self.serialiser). \
            with_retry_count(self.max_retries). \
            with_wire_type(self.wire_type).\
            with_timeout(2.0).build()
        rest_client = RestClient(client_config)
        self.assertEqual(rest_client.url, self.sample_url, "sample_urls do not match")
        self.assertEqual(rest_client.session.adapters["https://"].max_retries.total, self.max_retries)
        self.assertEqual(rest_client.session.adapters["http://"].max_retries.total, self.max_retries)
        self.assertEqual(type(rest_client.serde), self.serialiser.value, "serialiser is configured incorrectly")
        self.assertEqual(type(rest_client.wire), self.wire_type.value, "wire type is configured incorrectly")
        self.assertEqual(rest_client.timeout, 2.0, "timeout is configured incorrectly")

    def test_client_creation_failure(self):
        builder = RestClientConfigBuilder().with_url(self.sample_url)
        self.assertRaises(ValueError, builder.with_serialiser, "JSON")
        self.assertRaises(ValueError, builder.with_wire_type, "PROTOBUF")
        self.assertRaises(ValueError, builder.with_retry_count, "five")
        self.assertRaises(ValueError, builder.with_timeout, 0.005)

    @patch("rest.client.time.time_ns", return_value=_get_static_time())
    def test_get_stub_request(self, time_ns):
        rest_client = self._get_rest_client()
        ts = timestamp_pb2.Timestamp()
        ts.FromNanoseconds(time.time_ns())
        with patch("rest.client.uuid.uuid4", return_value=_get_static_uuid()):
            req = rest_client._get_stub_request()
            self.assertEqual(req.req_guid, _get_static_uuid())
            self.assertEqual(req.sent_time.seconds, ts.seconds)
            self.assertEqual(req.sent_time.nanos, ts.nanos)

    def test_uniqueness_of_stub_request(self):
        rest_client = self._get_rest_client()
        req1 = rest_client._get_stub_request()
        time.sleep(1)
        req2 = rest_client._get_stub_request()
        self.assertNotEqual(req1.req_guid, req2.req_guid)
        self.assertNotEqual(req1.sent_time.nanos, req2.sent_time.nanos)
        self.assertNotEqual(req1.sent_time.seconds, req2.sent_time.seconds)

    def test_client_send_success(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.return_value = self._get_stub_response()
        event_arr = [self._get_stub_event_payload()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = _get_static_uuid()
        req.req_guid = _get_static_uuid()
        time_in_ns = time.time_ns()
        req.sent_time.FromNanoseconds(time_in_ns)
        expected_req.sent_time.FromNanoseconds(time_in_ns)
        with patch("rest.client.requests.session", return_value=session_mock):
            rest_client = self._get_rest_client()
            expected_req.events.append(rest_client._convert_to_event_pb(self._get_stub_event_payload()))
            serialised_data = json_format.MessageToJson(expected_req)
            rest_client._get_stub_request = mock.MagicMock()
            rest_client._get_stub_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            rest_client.send(event_arr)
            post.assert_called_once_with(url=self.sample_url, data=serialised_data,
                                         headers={"Content-Type": "application/json"}, timeout=2.0)
            rest_client._parse_response.assert_called_once_with(post.return_value)

    def test_client_send_connection_failure(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.side_effect = ConnectionError("error connecting to host")
        event_arr = [self._get_stub_event_payload()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = _get_static_uuid()
        req.req_guid = _get_static_uuid()
        time_in_ns = time.time_ns()
        req.sent_time.FromNanoseconds(time_in_ns)
        expected_req.sent_time.FromNanoseconds(time_in_ns)
        with patch("rest.client.requests.session", return_value=session_mock):
            rest_client = self._get_rest_client()
            expected_req.events.append(rest_client._convert_to_event_pb(self._get_stub_event_payload()))
            serialised_data = json_format.MessageToJson(expected_req)
            rest_client._get_stub_request = mock.MagicMock()
            rest_client._get_stub_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            self.assertRaises(ConnectionError, rest_client.send, event_arr)
            post.assert_called_once_with(url=self.sample_url, data=serialised_data,
                                         headers={"Content-Type": "application/json"}, timeout=2.0)
            rest_client._parse_response.assert_not_called()

    def test_parse_response(self):
        resp = self._get_stub_response()
        rest_client = self._get_rest_client()
        deserialised_response = rest_client._parse_response(resp)
        self.assertEqual(deserialised_response.status, Status.STATUS_SUCCESS)
        self.assertEqual(deserialised_response.data["req_guid"], _get_static_uuid())
        self.assertEqual(deserialised_response.sent_time, _get_static_time())

    def _get_rest_client(self):
        client_config = RestClientConfigBuilder(). \
            with_url(self.sample_url). \
            with_serialiser(self.serialiser). \
            with_retry_count(self.max_retries). \
            with_timeout(2.0).build()
        return RestClient(client_config)

    def _get_stub_response(self):
        response = requests.Response()
        response.status_code = requests.status_codes.codes["ok"]
        json_response = {"status": 1, "code": 1, "sent_time": _get_static_time(),
                         "data": {"req_guid": _get_static_uuid()}}
        json_string2 = json.dumps(json_response)
        response._content = json_string2
        return response

    def _get_stub_event_payload(self):
        e = client.Event()
        e.type = "random_topic"
        e.event = {"a": "abc"}
        return e

    def _get_send_event_response(self):
        serialised_data = SendEventResponse()
