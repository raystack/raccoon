from serde.enum import Serialiser, WireType


class RestClientConfig:
    url: str
    max_retries: int
    serialiser: Serialiser
    headers: dict

    def __init__(self):
        self.headers = {}
        self.serialiser = Serialiser.JSON
        self.max_retries = 0
        self.wire_type = WireType.JSON


class RestClientConfigBuilder:

    def __init__(self):
        self.config = RestClientConfig()

    def with_url(self, url):
        self.config.url = url
        return self

    def with_retry_count(self, retry_count):
        if not isinstance(retry_count, int):
            raise ValueError("retry_count should be an integer")
        elif retry_count > 10:
            raise ValueError("retry should not be greater than 10")
        self.config.max_retries = retry_count
        return self

    def with_serialiser(self, content_type):
        if not isinstance(content_type, Serialiser):
            raise ValueError("invalid  serialiser/deserialiser type")
        self.config.serialiser = content_type
        return self

    def with_headers(self, headers):
        self.config.headers = headers

    def with_wire_type(self, wire_type):
        if not isinstance(wire_type, WireType):
            raise ValueError("invalid  serialiser/deserialiser type")
        self.config.wire_type = wire_type
        return self

    def build(self):
        return self.config
