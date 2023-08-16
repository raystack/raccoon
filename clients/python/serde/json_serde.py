import json

from google.protobuf import json_format


from serde.serde import Serde
from serde.wire import Wire


class JsonSerde(Serde, Wire):
    def serialise(self, event):
        return bytes(json.dumps(event), "utf-8")

    def get_content_type(self):
        return "application/json"

    def marshal(self, event):
        return json_format.MessageToJson(event)

    def unmarshal(self, data, template):
        return json_format.Parse(data, template)


