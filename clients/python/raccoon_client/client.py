from dataclasses import dataclass

from raccoon_client.protos.raystack.raccoon.v1beta1.raccoon_pb2 import (
    SendEventResponse,
)


@dataclass
class RaccoonResponseError(IOError):
    def __init__(self, status_code, msg):
        super().__init__(msg)
        self.status_code = status_code


@dataclass
class Event:
    type: str
    event: object


class Client:  # pylint: disable=too-few-public-methods
    def send(self, events: [Event]) -> (str, SendEventResponse, RaccoonResponseError):
        raise NotImplementedError()
