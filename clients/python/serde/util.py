from serde.enum import ContentType
from serde.json_serde import JsonSerde


def get_serde(content_type):
    if content_type == ContentType.JSON:
        return JsonSerde()
    else:
        return NotImplementedError()


CONTENT_TYPE_HEADER_KEY = "Content-Type"
