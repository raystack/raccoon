from enum import Enum

from serde.json_serde import JsonSerde


class Serialiser(Enum):
    JSON = JsonSerde
    PROTOBUF = 2


class WireType(Enum):
    JSON = JsonSerde
    PROTOBUF = 2
