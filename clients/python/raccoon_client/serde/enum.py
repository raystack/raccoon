from enum import Enum

from raccoon_client.serde.json_serde import JsonSerde
from raccoon_client.serde.protobuf_serde import ProtobufSerde


class Serialiser(Enum):
    JSON = JsonSerde
    PROTOBUF = ProtobufSerde


class WireType(Enum):
    JSON = JsonSerde
    PROTOBUF = ProtobufSerde
