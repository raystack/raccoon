import json

from google.protobuf import json_format
from google.protobuf.message import Message

from raccoon_client.protos.raystack.raccoon.v1beta1.raccoon_pb2 import (
    SendEventRequest,
    SendEventResponse,
)
from raccoon_client.serde.serde import Serde
from raccoon_client.serde.wire import Wire


class JsonSerde(Serde, Wire):
    # uses json.dumps since the input can be either protobuf message or dictionary
    def serialise(self, event):
        if isinstance(event, Message):
            return bytes(json_format.MessageToJson(event), "utf-8")
        return bytes(json.dumps(event), "utf-8")

    def get_content_type(self):
        return "application/json"

    def marshal(self, event: SendEventRequest):
        req_dict = json_format.MessageToDict(event, preserving_proto_field_name=True)
        req_dict["sent_time"] = {
            "seconds": event.sent_time.seconds,
            "nanos": event.sent_time.nanos,
        }
        return json.dumps(
            req_dict
        )  # uses json_format since the event is always a protobuf message

    def unmarshal(self, data, template: SendEventResponse):
        return json_format.Parse(data, template)
