from raccoon_client.serde.enum import Serialiser, WireType
from raccoon_client.serde.json_serde import JsonSerde
from raccoon_client.serde.protobuf_serde import ProtobufSerde
from raccoon_client.serde.serde import Serde
from raccoon_client.serde.wire import Wire


def get_serde(serialiser) -> Serde:
    if serialiser == Serialiser.JSON:
        return JsonSerde()
    if serialiser == Serialiser.PROTOBUF:
        return ProtobufSerde()
    raise ValueError()


def get_wire_type(wire_type) -> Wire:
    if wire_type == WireType.JSON:
        return JsonSerde()
    if wire_type == WireType.PROTOBUF:
        return ProtobufSerde()
    raise ValueError()


CONTENT_TYPE_HEADER_KEY = "Content-Type"
