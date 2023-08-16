import time
import uuid

import requests

from client import Client, Event
from protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest, SendEventResponse
from rest.option import RestClientConfig
from serde.util import get_serde, CONTENT_TYPE_HEADER_KEY


class RestClient(Client):
    session: requests.Session

    def __init__(self, config: RestClientConfig):
        self.session = requests.session()
        self.url = config.url
        self.serde = get_serde(config.content_type)
        self.headers = self._set_content_type_header(config.headers)
        self.max_retries = config.max_retries

    def send(self, events: [Event]):
        req = self._get_stub_request()
        for e in events:
            req.events.append(e)
        response = self.session.post(url=self.url, data=self.serde.serialise(req), headers=self.headers)
        deserialised_response = self._parse_response(response)
        return req.req_guid, deserialised_response, response

    def _get_stub_request(self):
        req = SendEventRequest()
        req.req_guid = uuid.uuid4()
        req.sent_time.FromNanoseconds(time.time_ns())
        return req

    def _set_content_type_header(self, headers):
        headers[CONTENT_TYPE_HEADER_KEY] = self.serde.get_content_type()
        return headers

    def _parse_response(self, response) -> SendEventResponse:
        if len(response.content) != 0:
            event_response = self.serde.deserialise(str(response.content), SendEventResponse())
            return event_response
