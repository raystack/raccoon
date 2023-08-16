import json
import time
import unittest

from unittest import mock
from unittest.mock import patch

import requests
from google.protobuf import json_format

import client
from protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest, SendEventResponse, Status
from rest.client import RestClient
from rest.option import RestClientConfigBuilder
from serde.enum import Serialiser


class RestClientTest(unittest.TestCase):
    sample_url = "http://localhost:8080/api/v1/"
    max_retries = 3
    content_type = Serialiser.JSON

    def test_client_creation(self):
        client_config = RestClientConfigBuilder(). \
            with_url(self.sample_url). \
            with_serialiser(self.content_type). \
            with_retry_count(self.max_retries).build()
        rest_client = RestClient(client_config)
        self.assertEqual(rest_client.url, self.sample_url, "sample_urls do not match")
        self.assertEqual(rest_client.session.adapters["https://"].max_retries.total, self.max_retries)
        self.assertEqual(rest_client.session.adapters["http://"].max_retries.total, self.max_retries)
        self.assertEqual(type(rest_client.serde), self.content_type.value, "content type is configured incorrectly")

    def test_get_stub_request(self):
        client = self._get_rest_client()

        with patch("rest.client.uuid.uuid4", return_value=self._get_static_uuid()):
            req = client._get_stub_request()
            self.assertEqual(req.req_guid, self._get_static_uuid())

    def test_client_send(self):
        session_mock = mock.Mock()
        post = mock.MagicMock()
        session_mock.post = post
        post.return_value = self._get_stub_response()
        event_arr = [self._get_stub_event_payload()]
        req = SendEventRequest()
        expected_req = SendEventRequest()
        expected_req.req_guid = self._get_static_uuid()
        req.req_guid = self._get_static_uuid()
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
                                         headers={"Content-Type": "application/json"})
            rest_client._parse_response.assert_called_once_with(post.return_value)

    def test_parse_response(self):
        resp = self._get_stub_response()
        rest_client = self._get_rest_client()
        deserialised_response = rest_client._parse_response(resp)
        self.assertEqual(deserialised_response.status, Status.STATUS_SUCCESS)
        self.assertEqual(deserialised_response.data["req_guid"], self._get_static_uuid())
        self.assertEqual(deserialised_response.sent_time, self._get_static_time())

    def _get_rest_client(self):
        client_config = RestClientConfigBuilder(). \
            with_url(self.sample_url). \
            with_serialiser(self.content_type). \
            with_retry_count(self.max_retries).build()
        return RestClient(client_config)

    def _get_stub_response(self):
        response = requests.Response()
        response.status_code = requests.status_codes.codes["ok"]
        json_response = {"status": 1, "code": 1, "sent_time": self._get_static_time(),
                         "data": {"req_guid": self._get_static_uuid()}}
        json_string2 = json.dumps(json_response)
        response._content = json_string2
        return response

    def _get_stub_event_payload(self):
        e = client.Event()
        e.type = "random_topic"
        e.event = {"a": "abc"}
        return e

    def _get_static_uuid(self):
        return "17e2ac19-df8b-4a30-b111-fd7f5073d2f5"

    def _get_static_time(self):
        return 1692182259

    def _get_send_event_response(self):
        serialised_data = SendEventResponse()
