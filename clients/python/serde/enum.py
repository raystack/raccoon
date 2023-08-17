from enum import Enum

from serde.json_serde import JsonSerde
from serde.protobuf_serde import ProtobufSerde


class Serialiser(Enum):
    JSON = JsonSerde
    PROTOBUF = ProtobufSerde


class WireType(Enum):
    JSON = JsonSerde
    PROTOBUF = ProtobufSerde
