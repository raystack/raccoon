from google.protobuf.message import Message

from serde.serde import Serde
from serde.wire import Wire


class ProtobufSerde(Serde, Wire):
    def serialise(self, event: Message):
        return event.SerializeToString()  # the name is a misnomer, returns bytes

    def marshal(self, obj: Message):
        return obj.SerializeToString()

    def unmarshal(self, marshalled_data: bytes, template: Message):
        template.ParseFromString(marshalled_data)
        return template

    def get_content_type(self):
        return "application/proto"
