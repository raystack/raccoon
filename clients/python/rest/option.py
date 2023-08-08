from clients.python.serde.enum import ContentType


class Config:
    url: str
    max_retries: int
    content_type: ContentType

