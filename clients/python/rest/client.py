import time
import uuid

import requests

from clients.python.client import Client, Event
from clients.python.protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest
from clients.python.serde.util import get_serde, CONTENT_TYPE_HEADER_KEY
from google.protobuf import timestamp_pb2


class RestClient(Client):
    session: requests.Session
    HTTP_PATH = "/api/v1/events"

    def __init__(self, config):
        self.session = requests.session()
        self.url = config.url + self.HTTP_PATH
        self.serde = get_serde(config.content_type)
        self.headers = self._set_content_type_header({})

    def send(self, events: [Event]):
        req = self._get_stub_request()
        req.events = events
        self.session.post(url=self.url, data=self.serde.serialise(req), headers=self.headers)

    def _get_stub_request(self):
        req = SendEventRequest()
        req.req_guid = uuid.uuid4()
        req.sent_time = timestamp_pb2.Timestamp.FromNanoseconds(time.time_ns())
        return req

    def _set_content_type_header(self, headers):
        headers[CONTENT_TYPE_HEADER_KEY] = self.serde.get_content_type()
        return headers

