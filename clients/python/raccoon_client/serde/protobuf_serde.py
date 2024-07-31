from google.protobuf.message import Message

from raccoon_client.serde.serde import Serde
from raccoon_client.serde.wire import Wire


class ProtobufSerde(Serde, Wire):
    def serialise(self, event: Message):
        if not isinstance(event, Message):
            raise ValueError("event should be a protobuf message")
        return event.SerializeToString()  # the name is a misnomer, returns bytes

    def marshal(self, event: Message):
        return event.SerializeToString()

    def unmarshal(self, data: bytes, template: Message):
        template.ParseFromString(data)
        return template

    def get_content_type(self):
        return "application/proto"
