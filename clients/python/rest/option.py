from serde.enum import ContentType


class RestClientConfig:
    url: str
    max_retries: int
    content_type: ContentType

