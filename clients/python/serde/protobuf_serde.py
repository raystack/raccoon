from google._upb._message import Message

from serde.serde import Serde
from serde.wire import Wire


class ProtobufSerde(Serde, Wire):
    def serialise(self, event: Message):
        return bytes(event.SerializeToString(), "utf-8")

    def marshal(self, obj: Message):
        return bytes(obj.SerializeToString(), "utf-8")

    def unmarshal(self, serialised_data: bytes, template: Message):
        return template.ParseFromString(serialised_data)

    def get_content_type(self):
        return "application/proto"
