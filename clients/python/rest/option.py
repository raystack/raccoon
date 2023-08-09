from serde.enum import ContentType


class RestClientConfig:
    url: str
    max_retries: int
    content_type: ContentType
    headers: dict

    def __init__(self):
        self.headers = {}
        self.content_type = ContentType.JSON
        self.max_retries = 0


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

    def with_content_type(self, content_type):
        if not isinstance(content_type, ContentType):
            raise ValueError("invalid  serialiser/deserialiser type")
        self.config.content_type = content_type
        return self

    def with_headers(self, headers):
        self.config.headers = headers

    def build(self):
        return self.config
