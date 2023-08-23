from raccoon_client.serde.enum import Serialiser, WireType
from raccoon_client.serde.json_serde import JsonSerde
from raccoon_client.serde.protobuf_serde import ProtobufSerde
from raccoon_client.serde.serde import Serde
from raccoon_client.serde.wire import Wire


def get_serde(serialiser) -> Serde:
    if serialiser == Serialiser.JSON:
        return JsonSerde()
    elif serialiser == Serialiser.PROTOBUF:
        return ProtobufSerde()
    else:
        raise ValueError()


def get_wire_type(wire_type) -> Wire:
    if wire_type == WireType.JSON:
        return JsonSerde()
    elif wire_type == WireType.PROTOBUF:
        return ProtobufSerde()
    else:
        raise ValueError()


CONTENT_TYPE_HEADER_KEY = "Content-Type"
