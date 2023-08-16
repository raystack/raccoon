from serde.enum import Serialiser, WireType
from serde.json_serde import JsonSerde
from serde.serde import Serde
from serde.wire import Wire


def get_serde(serialiser) -> Serde:
    if serialiser == Serialiser.JSON:
        return JsonSerde()
    else:
        raise NotImplementedError()

def get_wire_type(wire_type) -> Wire:
    if wire_type == WireType.JSON:
        return JsonSerde()
    else:
        raise NotImplementedError()


CONTENT_TYPE_HEADER_KEY = "Content-Type"