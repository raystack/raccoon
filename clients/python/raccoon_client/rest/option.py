from dataclasses import dataclass

from raccoon_client.serde.enum import Serialiser, WireType


@dataclass
class HttpConfig:
    url: str
    max_retries: int
    timeout: float
    headers: dict[str]

    def __init__(self):
        self.url = ""
        self.max_retries = 3
        self.timeout = 1.0
        self.headers = {}

    def clone(self):
        cloned_config = HttpConfig()
        cloned_config.url = self.url
        cloned_config.max_retries = self.max_retries
        cloned_config.timeout = self.timeout
        cloned_config.headers = self.headers
        return cloned_config


@dataclass
class RestClientConfig:
    serialiser: Serialiser
    wire_type: WireType
    http: HttpConfig

    def __init__(self):
        self.serialiser = Serialiser.JSON
        self.wire_type = WireType.JSON
        self.http = HttpConfig()

    def get_http_config(self):
        return self.http.clone()


class RestClientConfigBuilder:
    def __init__(self):
        self.config = RestClientConfig()

    def with_url(self, url):
        self.config.http.url = url
        return self

    def with_retry_count(self, retry_count):
        if not isinstance(retry_count, int):
            raise ValueError("retry_count should be an integer")
        if retry_count > 10:
            raise ValueError("retry should not be greater than 10")
        self.config.http.max_retries = retry_count
        return self

    def with_serialiser(self, content_type):
        if not isinstance(content_type, Serialiser):
            raise ValueError("invalid  serialiser/deserialiser type")
        self.config.serialiser = content_type
        return self

    def with_headers(self, headers):
        self.config.http.headers = headers

    def with_wire_type(self, wire_type):
        if not isinstance(wire_type, WireType):
            raise ValueError("invalid  serialiser/deserialiser type")
        self.config.wire_type = wire_type
        return self

    def with_timeout(self, timeout):
        if not isinstance(timeout, float):
            raise ValueError
        if timeout > 10:
            raise ValueError("timeout too high")
        if timeout < 0.010:
            raise ValueError("timeout is too low")
        self.config.http.timeout = timeout
        return self

    def build(self):
        return self.config
