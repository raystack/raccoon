import json

from clients.python.serde.serde import Serde


class JsonSerde(Serde):
    def serialise(self, event):
        json.dumps(event)

    def deserialise(self, data):
        json.loads(data)

    def get_content_type(self):
        return "application/json"


