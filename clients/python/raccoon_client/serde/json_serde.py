import json

from google.protobuf import json_format
from google.protobuf.message import Message

from raccoon_client.serde.serde import Serde
from raccoon_client.serde.wire import Wire


class JsonSerde(Serde, Wire):
    def serialise(self, event):
        if isinstance(event, Message):
            return bytes(json_format.MessageToJson(event), "utf-8")
        return bytes(json.dumps(event), "utf-8")  # uses json.dumps since the input can be either protobuf message or dictionary

    def get_content_type(self):
        return "application/json"

    def marshal(self, event):
        return json_format.MessageToJson(event)  # uses json_format since the event is always a protobuf message

    def unmarshal(self, data, template):
        return json_format.Parse(data, template)


