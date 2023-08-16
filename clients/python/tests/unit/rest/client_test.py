import json
import time
import unittest
import uuid
from unittest import mock
from unittest.mock import patch

import requests
from google.protobuf import json_format

from protos.raystack.raccoon.v1beta1.raccoon_pb2 import Event, SendEventRequest, SendEventResponse
from rest.client import RestClient
from rest.option import RestClientConfigBuilder
from serde.enum import ContentType


class RestClientTest(unittest.TestCase):

    sample_url = "http://localhost:8080/api/v1/"
    max_retries = 3
    content_type = ContentType.JSON

    def test_client_creation(self):
        client_config = RestClientConfigBuilder().\
            with_url(self.sample_url).\
            with_content_type(self.content_type).\
            with_retry_count(self.max_retries).build()
        rest_client = RestClient(client_config)
        self.assertEqual(rest_client.url, self.sample_url, "sample_urls do not match")
        self.assertEqual(rest_client.max_retries, self.max_retries, "max_retries is configured incorrectly")
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
        expected_req.events.append(self._get_stub_event_payload())
        serialised_data = json_format.MessageToJson(expected_req)
        with patch("rest.client.requests.session", return_value=session_mock):
            rest_client = self._get_rest_client()
            rest_client._get_stub_request = mock.MagicMock()
            rest_client._get_stub_request.return_value = req
            rest_client._parse_response = mock.MagicMock()
            rest_client.send(event_arr)
            post.assert_called_once_with(url=self.sample_url, data=serialised_data, headers={"Content-Type": "application/json"})
            rest_client._parse_response.assert_called_once_with(post.return_value)

    def _get_rest_client(self):
        client_config = RestClientConfigBuilder().\
            with_url(self.sample_url).\
            with_content_type(self.content_type).\
            with_retry_count(self.max_retries).build()
        return RestClient(client_config)

    def _get_stub_response(self):
        response = requests.Response()
        response.status_code = requests.status_codes.codes["ok"]
        response._content = bytes("some_raw_data", "utf-8")
        return response

    def _get_stub_event_payload(self):
        e = Event()
        e.type = "random_topic"
        e.event_bytes = bytes("random_bytes", "utf-8")
        return e

    def _get_static_uuid(self):
        return "17e2ac19-df8b-4a30-b111-fd7f5073d2f5"

    def _get_send_event_response(self):
        serialised_data = SendEventResponse()

