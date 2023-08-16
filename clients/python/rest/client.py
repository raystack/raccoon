import time
import uuid

import requests

from client import Client, Event
from protos.raystack.raccoon.v1beta1.raccoon_pb2 import SendEventRequest, SendEventResponse, Event as EventPb
from rest.option import RestClientConfig
from serde.serde import Serde
from serde.util import get_serde, CONTENT_TYPE_HEADER_KEY, get_wire_type
from serde.wire import Wire


class RestClient(Client):
    session: requests.Session
    serde: Serde
    wire: Wire

    def __init__(self, config: RestClientConfig):
        self.session = requests.session()
        self.url = config.url
        self.serde = get_serde(config.serialiser)
        self.wire = get_wire_type(config.wire_type)
        self.headers = self._set_content_type_header(config.headers)
        self.max_retries = config.max_retries

    def send(self, events: [Event]):
        req = self._get_stub_request()
        events_pb = map(lambda x: self._convert_to_event_pb(x), events)
        req.events.extend(events_pb)
        response = self.session.post(url=self.url, data=self.wire.marshal(req), headers=self.headers)
        deserialised_response = self._parse_response(response)
        return req.req_guid, deserialised_response, response

    def _convert_to_event_pb(self, e: Event):
        proto_event = EventPb()
        proto_event.event_bytes = self.serde.serialise(e.event)
        proto_event.type = e.type
        return proto_event

    def _get_stub_request(self):
        req = SendEventRequest()
        req.req_guid = uuid.uuid4()
        req.sent_time.FromNanoseconds(time.time_ns())
        return req

    def _set_content_type_header(self, headers):
        headers[CONTENT_TYPE_HEADER_KEY] = self.wire.get_content_type()
        return headers

    def _parse_response(self, response) -> SendEventResponse:
        if len(response.content) != 0:
            event_response = self.wire.unmarshal(str(response.content), SendEventResponse())
            return event_response
