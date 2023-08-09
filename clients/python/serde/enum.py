from enum import Enum

from serde.json_serde import JsonSerde


class ContentType(Enum):
    JSON = JsonSerde
    PROTOBUF = 2



